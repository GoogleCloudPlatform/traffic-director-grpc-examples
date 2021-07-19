/*
 *
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"

	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	walletpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet"
	accountpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/account"
	statspb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/stats"
	"google.golang.org/grpc/grpc-wallet/observability"
	"google.golang.org/grpc/grpc-wallet/utility"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/xds"
)

var (
	port             = flag.String("port", "18881", "the port to listen on, default '18881'")
	adminPort        = flag.String("admin_port", "28881", "the port to listen on, for admin services like CSDS, health, channelz etc, default '28881'")
	statsServer      = flag.String("stats_server", "localhost:18882", "address of the stats service, default 'localhost:18882'")
	accountServer    = flag.String("account_server", "localhost:18883", "address of the account service, default 'localhost:18883'")
	hostnameSuffix   = flag.String("hostname_suffix", "", "suffix to append to hostname in response header for outgoing RPCs, default ''")
	gcpClientProject = flag.String("gcp_client_project", "", "if set, metrics and traces will be sent to Cloud Monitoring and Cloud Trace")
	creds            = flag.String("creds", "insecure", "type of transport credentials to use. Supported values include 'xds' and 'insecure', defaults to 'insecure'")
	v1Behavior       = flag.Bool("v1_behavior", false, "If true, only aggregate balance is reported. Default is 'false'")
)

var users = map[string]map[string]int{
	"Alice": {
		"cd0aa985": 314,
		"454349e4": 159,
	},
	"Bob": {
		"148de9c5": 271,
		"2e7d2c03": 828,
	},
}

// server provides an implementation of the Wallet service as defined in
// proto/grpc/examples/wallet/wallet.proto.
type server struct {
	walletpb.UnimplementedWalletServer
	v1Behavior    bool
	hostname      string
	accountClient accountpb.AccountClient
	statsClient   statspb.StatsClient
}

// buildBalanceResponse is a helper method which uses the price from the stats
// server to build a BalanceResponse.
func buildBalanceResponse(v1Behavior bool, name string, price int64, req *walletpb.BalanceRequest) (*walletpb.BalanceResponse, error) {
	result := walletpb.BalanceResponse{}
	addresses, ok := users[name]
	if !ok {
		return &result, status.Errorf(codes.NotFound, "could not identify user: %v", name)
	}
	total := int64(0)
	for address, count := range addresses {
		balance := int64(count) * price
		total += balance
		if !v1Behavior && req.GetIncludeBalancePerAddress() {
			result.Addresses = append(result.Addresses, &walletpb.BalancePerAddress{
				Address: address,
				Balance: balance,
			})
		}
	}
	result.Balance = total
	return &result, nil
}

// FetchBalance implements the FetchBalance method from the Wallet service.
func (s *server) FetchBalance(ctx context.Context, req *walletpb.BalanceRequest) (*walletpb.BalanceResponse, error) {
	// Validate membership using incoming metadata.
	_, name, token, membership, err := utility.ValidateMembership(ctx, s.accountClient)
	if err != nil {
		return nil, err
	}

	// Add the hostname header.
	if err := grpc.SetHeader(ctx, metadata.Pairs("hostname", s.hostname)); err != nil {
		return nil, err
	}

	// Get the price.
	md := metadata.Pairs("authorization", token, "membership", membership)
	statsCtx := metadata.NewOutgoingContext(ctx, md)
	var header metadata.MD
	r, err := s.statsClient.FetchPrice(statsCtx, &statspb.PriceRequest{}, grpc.Header(&header))
	utility.PrintHostname(header)
	log.Printf("grpc-coin price: %d", r.GetPrice())

	// Calculate the result.
	return buildBalanceResponse(s.v1Behavior, name, r.GetPrice(), req)
}

// WatchBalance implements the WatchBalance method from the Wallet service.
func (s *server) WatchBalance(req *walletpb.BalanceRequest, srv walletpb.Wallet_WatchBalanceServer) error {
	// Validate membership using incoming metadata.
	accountCtx := srv.Context()
	_, name, token, membership, err := utility.ValidateMembership(accountCtx, s.accountClient)
	if err != nil {
		return err
	}

	// Add the hostname header.
	if err := srv.SetHeader(metadata.Pairs("hostname", s.hostname)); err != nil {
		return err
	}

	// Connect to the stats client.
	md := metadata.Pairs("authorization", token, "membership", membership)
	statsCtx := metadata.NewOutgoingContext(srv.Context(), md)
	statsSrv, err := s.statsClient.WatchPrice(statsCtx, &statspb.PriceRequest{})
	if err != nil {
		return status.Errorf(codes.Unavailable, "error establishing stats.WatchPrice stream: %v", err)
	}
	header, err := statsSrv.Header()
	if err != nil {
		return status.Errorf(codes.Internal, "could not extract header: %v", err)
	}
	utility.PrintHostname(header)

	for {
		// Receive the price.
		r, err := statsSrv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		price := r.GetPrice()
		log.Printf("grpc-coin price: %d", price)

		// Send balance.
		result, err := buildBalanceResponse(s.v1Behavior, name, price, req)
		if err != nil {
			return err
		}
		if err := srv.Send(result); err != nil {
			return err
		}
	}
}

// grpcServer wraps methods that are invoked on a vanilla gRPC server or an
// xds-enabled gRPC server.
type grpcServer interface {
	grpc.ServiceRegistrar
	reflection.GRPCServer
	Serve(net.Listener) error
}

func main() {
	flag.Parse()

	// Parse credentials type from the command line to determine if xDS
	// credentials are to be used.
	xdsCreds, err := utility.ParseCredentialsType(*creds)
	if err != nil {
		log.Fatal(err)
	}

	// Use insecure credentials by default. But if xDS credentials are specified
	// on the command line, we create client and server xDS credentials with an
	// insecure fallback credentials.
	clientCreds, serverCreds := insecure.NewCredentials(), insecure.NewCredentials()
	if xdsCreds {
		var err error
		clientCreds, err = xdscreds.NewClientCredentials(xdscreds.ClientOptions{FallbackCreds: insecure.NewCredentials()})
		if err != nil {
			log.Fatalf("Failed to create client xDS credentials: %v", err)
		}
		serverCreds, err = xdscreds.NewServerCredentials(xdscreds.ServerOptions{FallbackCreds: insecure.NewCredentials()})
		if err != nil {
			log.Fatalf("Failed to create server xDS credentials: %v", err)
		}
	}

	// Create client and server transport credentials as dial options and server
	// options respectively.
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(clientCreds)}
	serverOpts := []grpc.ServerOption{grpc.Creds(serverCreds)}

	// Export stats to Stackdriver if a GCP project was specified.
	if *gcpClientProject != "" {
		sd := observability.ConfigureStackdriver(*gcpClientProject)
		defer sd.Flush()
		defer sd.StopMetricsExporter()
		dialOpts = append(dialOpts, grpc.WithStatsHandler(new(ocgrpc.ClientHandler)))
		serverOpts = append(serverOpts, grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	}

	// Dial the account server.
	accountConn, err := grpc.Dial(*accountServer, dialOpts...)
	if err != nil {
		log.Fatalf("Failed to dail the account server at %q: %v", *accountServer, err)
	}
	defer accountConn.Close()
	accountClient := accountpb.NewAccountClient(accountConn)

	// Dial the stats server.
	statsConn, err := grpc.Dial(*statsServer, dialOpts...)
	if err != nil {
		log.Fatalf("Failed to dail the stats server at %q: %v", *statsServer, err)
	}
	defer statsConn.Close()
	statsClient := statspb.NewStatsClient(statsConn)

	// Create a listener to serve Wallet service RPCs on -port.
	walletListener, err := net.Listen("tcp4", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create the Wallet service implementation.
	walletServer := &server{
		v1Behavior:    *v1Behavior,
		hostname:      utility.GenHostname(*hostnameSuffix),
		accountClient: accountClient,
		statsClient:   statsClient,
	}

	// Start serving Wallet service on -port. We use an xDS-enabled gRPC server
	// when xDS credentials are to be used, and a vanilla gRPC server otherwise.
	var s grpcServer
	if xdsCreds {
		s = xds.NewGRPCServer(serverOpts...)
	} else {
		s = grpc.NewServer(serverOpts...)
	}
	walletpb.RegisterWalletServer(s, walletServer)

	// Also export the gRPC health check service on the same port as that of the
	// Wallet service. This will make it easier to configure GCP health checks
	// with the --use-serving-port flag.
	//
	// Note that the GCP health check will not work if an xds-enabled gRPC server
	// is used in mTLS mode. In this case, the health check on the -admin_port can
	// be used.
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, healthServer)
	reflection.Register(s)
	go func() {
		if err := s.Serve(walletListener); err != nil {
			log.Fatalf("Failed to serve Wallet service: %v", err)
		}
	}()

	// Expose admin services (csds, health and reflection) on -admin_port.
	if err := utility.StartAdminServices(*adminPort); err != nil {
		log.Fatalf("StartAdminServices(%s) failed: %v", *adminPort, err)
	}
}
