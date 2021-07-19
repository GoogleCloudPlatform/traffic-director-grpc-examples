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
	"log"
	"math"
	"net"
	"time"

	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
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
	port             = flag.String("port", "18882", "the port to listen on, default '18882'")
	adminPort        = flag.String("admin_port", "28882", "the port to listen on, for admin services like CSDS, health, channelz etc, default '28882'")
	accountServer    = flag.String("account_server", "localhost:18883", "address of the account service, default 'localhost:18883'")
	hostnameSuffix   = flag.String("hostname_suffix", "", "suffix to append to hostname in response header for outgoing RPCs, default ''")
	gcpClientProject = flag.String("gcp_client_project", "", "if set, metrics and traces will be sent to Cloud Monitoring and Cloud Trace")
	creds            = flag.String("creds", "insecure", "type of transport credentials to use. Supported values include 'xds' and 'insecure', defaults to 'insecure'")
	premiumOnly      = flag.Bool("premium_only", false, "whether this service is for users with premium access only, default 'false'")
)

// server provides an implementation of the Stats service as defined in
// proto/grpc/examples/wallet/stats/stats.proto.
type server struct {
	statspb.UnimplementedStatsServer
	premiumOnly   bool
	hostname      string
	accountClient accountpb.AccountClient
}

// getPrice calculates the price with the equation:
// price = 1000sin(time * 1/173) + 10000.
func getPrice() int64 {
	return int64(math.Sin(float64(time.Now().UnixNano()/int64(time.Millisecond))/173)*1000 + 10000)
}

// FetchPrice implements the FetchPrice method from the Stats service.
func (s *server) FetchPrice(ctx context.Context, req *statspb.PriceRequest) (*statspb.PriceResponse, error) {
	// Validate membership using incoming metadata.
	premium, _, _, _, err := utility.ValidateMembership(ctx, s.accountClient)
	if err != nil {
		return nil, err
	}

	// Reject if non-premium when premium only.
	if s.premiumOnly && !premium {
		return nil, status.Error(codes.PermissionDenied, "non-premium request while premium-only")
	}

	// Add the hostname header.
	if err := grpc.SetHeader(ctx, metadata.Pairs("hostname", s.hostname)); err != nil {
		return nil, err
	}

	return &statspb.PriceResponse{Price: getPrice()}, nil
}

// WatchPrice implements the WatchPrice method from the Stats service.
func (s *server) WatchPrice(req *statspb.PriceRequest, srv statspb.Stats_WatchPriceServer) error {
	// Validate membership using incoming metadata.
	ctx := srv.Context()
	premium, _, _, _, err := utility.ValidateMembership(ctx, s.accountClient)
	if err != nil {
		return err
	}

	// Reject if non-premium when premium only.
	if s.premiumOnly && !premium {
		return status.Error(codes.PermissionDenied, "non-premium request while premium-only")
	}

	// Add the hostname header.
	if err := srv.SetHeader(metadata.Pairs("hostname", s.hostname)); err != nil {
		return err
	}

	// Send & sleep according to membership status.
	d := time.Second
	if premium {
		d = time.Millisecond * 100
	}
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			if err := srv.Send(&statspb.PriceResponse{Price: getPrice()}); err != nil {
				return err
			}
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

	// Create a listener to serve Stats service RPCs on -port.
	statsListener, err := net.Listen("tcp4", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create the Stats service implementation.
	statsServer := &server{
		premiumOnly:   *premiumOnly,
		hostname:      utility.GenHostname(*hostnameSuffix),
		accountClient: accountClient,
	}

	// Start serving Stats service on -port. We use an xDS-enabled gRPC server
	// when xDS credentials are to be used, and a vanilla gRPC server otherwise.
	var s grpcServer
	if xdsCreds {
		s = xds.NewGRPCServer(serverOpts...)
	} else {
		s = grpc.NewServer(serverOpts...)
	}
	statspb.RegisterStatsServer(s, statsServer)

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
		if err := s.Serve(statsListener); err != nil {
			log.Fatalf("Failed to serve Stats service: %v", err)
		}
	}()

	// Expose admin services (csds, health and reflection) on -admin_port.
	if err := utility.StartAdminServices(*adminPort); err != nil {
		log.Fatalf("StartAdminServices(%s) failed: %v", *adminPort, err)
	}
}
