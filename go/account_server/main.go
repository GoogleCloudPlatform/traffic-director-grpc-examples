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
	"net"

	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	accountpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/account"
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
	port             = flag.String("port", "18883", "the port to listen on, default '18883'")
	adminPort        = flag.String("admin_port", "28883", "the port to listen on, for admin services like CSDS, health, channelz etc, default '28883'")
	hostnameSuffix   = flag.String("hostname_suffix", "", "suffix to append to hostname in response header for outgoing RPCs, default ''")
	gcpClientProject = flag.String("gcp_client_project", "", "if set, metrics and traces will be sent to Cloud Monitoring and Cloud Trace")
	creds            = flag.String("creds", "insecure", "type of transport credentials to use. Supported values include 'xds' and 'insecure', defaults to 'insecure'")
)

type user struct {
	name       string
	membership accountpb.MembershipType
}

var users = map[string]user{
	"2bd806c9": {"Alice", accountpb.MembershipType_PREMIUM},
	"81b637d8": {"Bob", accountpb.MembershipType_NORMAL},
}

// server provides an implementation of the Account service as defined in
// proto/grpc/examples/wallet/account/account.proto.
type server struct {
	accountpb.UnimplementedAccountServer
	hostname string
}

// GetUserInfo implements the GetUserInfo method from the Account service.
func (s *server) GetUserInfo(ctx context.Context, req *accountpb.GetUserInfoRequest) (*accountpb.GetUserInfoResponse, error) {
	// Get the user matching the request's token.
	usr, ok := users[req.Token]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unrecognized user token")
	}

	// Add the hostname header.
	if err := grpc.SetHeader(ctx, metadata.Pairs("hostname", s.hostname)); err != nil {
		return nil, err
	}

	log.Printf("Received: '%v'. Sending: '%v'", req.Token, usr)
	return &accountpb.GetUserInfoResponse{Name: usr.name, Membership: usr.membership}, nil
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
	log.Println("Using xDS credentials:", xdsCreds)

	// Use insecure credentials by default. But if xDS credentials are specified
	// on the command line, use xDS credentials with an insecure fallback.
	serverCreds := insecure.NewCredentials()
	if xdsCreds {
		var err error
		serverCreds, err = xdscreds.NewServerCredentials(xdscreds.ServerOptions{FallbackCreds: insecure.NewCredentials()})
		if err != nil {
			log.Fatalf("Failed to create server xDS credentials: %v", err)
		}
	}
	serverOpts := []grpc.ServerOption{grpc.Creds(serverCreds)}

	// Export stats to Stackdriver if a GCP project was specified.
	if *gcpClientProject != "" {
		sd := observability.ConfigureStackdriver(*gcpClientProject)
		defer sd.Flush()
		defer sd.StopMetricsExporter()
		serverOpts = append(serverOpts, grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	}

	// Create a listener to serve Account service RPCs on -port.
	accountListener, err := net.Listen("tcp4", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create the Account service implementation.
	accountServer := &server{hostname: utility.GenHostname(*hostnameSuffix)}

	// Start serving Account service on -port. We use an xDS-enabled gRPC server
	// when xDS credentials are to be used, and a vanilla gRPC server otherwise.
	var s grpcServer
	if xdsCreds {
		s = xds.NewGRPCServer(serverOpts...)
	} else {
		s = grpc.NewServer(serverOpts...)
	}
	accountpb.RegisterAccountServer(s, accountServer)

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
		if err := s.Serve(accountListener); err != nil {
			log.Fatalf("Failed to serve Account service: %v", err)
		}
	}()

	// Expose admin services (csds, health and reflection) on -admin_port.
	if err := utility.StartAdminServices(*adminPort); err != nil {
		log.Fatalf("StartAdminServices(%s) failed: %v", *adminPort, err)
	}
}
