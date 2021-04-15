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
	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/codes"
	accountpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/account"
	statspb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/stats"
	"google.golang.org/grpc/grpc-wallet/observability"
	"google.golang.org/grpc/grpc-wallet/utility"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	_ "google.golang.org/grpc/xds" // To enable xds support.
)

type arguments struct {
	port                 string
	adminPort            string
	accountServer        string
	hostnameSuffix       string
	premiumOnly          bool
	observabilityProject string
}

// parseArguments parses the command line arguments using the flag package.
func parseArguments() arguments {
	result := arguments{}
	flag.StringVar(&result.port, "port", "18882", "the port to listen on, default '18882'")
	flag.StringVar(&result.adminPort, "admin_port", "58882", "the admin port to listen on, default '58882'")
	flag.StringVar(&result.accountServer, "account_server", "localhost:18883", "address of the account service, default 'localhost:18883'")
	flag.StringVar(&result.hostnameSuffix, "hostname_suffix", "", "suffix to append to hostname in response header for outgoing RPCs, default ''")
	flag.BoolVar(&result.premiumOnly, "premium_only", false, "whether this service is for users with premium access only, default false")
	flag.StringVar(&result.observabilityProject, "observability_project", "", "if set, metrics and traces will be sent to Cloud Monitoring and Cloud Trace")
	flag.Parse()
	return result
}

type server struct {
	statspb.UnimplementedStatsServer
	args          *arguments
	hostname      string
	accountClient accountpb.AccountClient
}

// getPrice calculates the price with the equation price = 1000sin(time * 1/173) + 10000.
func getPrice() int64 {
	return int64(math.Sin(float64(time.Now().UnixNano()/int64(time.Millisecond))/173)*1000 + 10000)
}

func (s *server) FetchPrice(ctx context.Context, req *statspb.PriceRequest) (*statspb.PriceResponse, error) {
	// Validate membership using incoming metadata.
	premium, _, _, _, err := utility.ValidateMembership(ctx, s.accountClient)
	if err != nil {
		return nil, err
	}

	// Reject if non-premium when premium only.
	if s.args.premiumOnly && !premium {
		return nil, status.Error(codes.PermissionDenied, "non-premium request while premium-only")
	}

	// Add the hostname header.
	if err := grpc.SetHeader(ctx, metadata.Pairs("hostname", s.hostname)); err != nil {
		return nil, err
	}

	return &statspb.PriceResponse{Price: getPrice()}, nil
}

func (s *server) WatchPrice(req *statspb.PriceRequest, srv statspb.Stats_WatchPriceServer) error {
	// Validate membership using incoming metadata.
	ctx := srv.Context()
	premium, _, _, _, err := utility.ValidateMembership(ctx, s.accountClient)
	if err != nil {
		return err
	}

	// Reject if non-premium when premium only.
	if s.args.premiumOnly && !premium {
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

func main() {
	args := parseArguments()

	var dialOpts = []grpc.DialOption{grpc.WithInsecure()}
	var serverOpts []grpc.ServerOption
	if args.observabilityProject != "" {
		sd := observability.ConfigureStackdriver(args.observabilityProject)
		defer sd.Flush()
		defer sd.StopMetricsExporter()
		dialOpts = append(dialOpts, grpc.WithStatsHandler(new(ocgrpc.ClientHandler)))
		serverOpts = append(serverOpts, grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	}

	// Dial account server.
	conn, err := grpc.Dial(args.accountServer, dialOpts...)
	if err != nil {
		log.Fatalf("did not connect: %v.", err)
	}
	defer conn.Close()
	c := accountpb.NewAccountClient(conn)

	// Start admin server
	adminListener, err := net.Listen("tcp", "localhost:"+args.adminPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	adminServer := grpc.NewServer()
	cleanup, err := admin.Register(adminServer)
	if err != nil {
		log.Fatalf("failed to register admin: %v", err)
	}
	defer cleanup()
	go adminServer.Serve(adminListener)
	defer adminServer.Stop()

	// Listen & serve.
	lis, err := net.Listen("tcp", ":"+args.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(serverOpts...)
	statspb.RegisterStatsServer(s, &server{
		args:          &args,
		hostname:      utility.GenHostname(args.hostnameSuffix),
		accountClient: c,
	})
	reflection.Register(s)
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, healthServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
