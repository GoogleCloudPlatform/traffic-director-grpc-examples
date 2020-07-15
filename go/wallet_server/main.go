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
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	walletpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet"
	accountpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/account"
	statspb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/stats"
	"google.golang.org/grpc/grpc-wallet/utility"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	_ "google.golang.org/grpc/xds" // To enable xds support.
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

type arguments struct {
	port           string
	accountServer  string
	statsServer    string
	hostnameSuffix string
	v1Behavior     bool
}

// parseArguments parses the command line arguments using the flag package.
func parseArguments() arguments {
	result := arguments{}
	flag.StringVar(&result.port, "port", "18881", "the port to listen on, default '18881'")
	flag.StringVar(&result.accountServer, "account_server", "localhost:18883", "address of the account service, default 'localhost:18883'")
	flag.StringVar(&result.statsServer, "stats_server", "localhost:18882", "address of the stats service, default 'localhost:18882'")
	flag.StringVar(&result.hostnameSuffix, "hostname_suffix", "", "suffix to append to hostname in response header for outgoing RPCs, default ''")
	flag.BoolVar(&result.v1Behavior, "v1_behavior", false, "usage")
	flag.Parse()
	return result
}

type server struct {
	walletpb.UnimplementedWalletServer
	args          *arguments
	hostname      string
	accountClient accountpb.AccountClient
	statsClient   statspb.StatsClient
}

// buildBalanceResponse uses the price from the stats server to build a BalanceResponse.
func buildBalanceResponse(args *arguments, name string, price int64, req *walletpb.BalanceRequest) (*walletpb.BalanceResponse, error) {
	result := walletpb.BalanceResponse{}
	addresses, ok := users[name]
	if !ok {
		return &result, status.Errorf(codes.NotFound, "could not identify user: %v", name)
	}
	total := int64(0)
	for address, count := range addresses {
		balance := int64(count) * price
		total += balance
		if !args.v1Behavior && req.GetIncludeBalancePerAddress() {
			result.Addresses = append(result.Addresses, &walletpb.BalancePerAddress{
				Address: address,
				Balance: balance,
			})
		}
	}
	result.Balance = total
	return &result, nil
}

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
	statsCtx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD
	r, err := s.statsClient.FetchPrice(statsCtx, &statspb.PriceRequest{}, grpc.Header(&header))
	utility.PrintHostname(header)
	log.Printf("grpc-coin price: %d", r.GetPrice())

	// Calculate the result.
	return buildBalanceResponse(s.args, name, r.GetPrice(), req)
}

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
	statsCtx := metadata.NewOutgoingContext(context.Background(), md)
	statsSrv, err := s.statsClient.WatchPrice(statsCtx, &statspb.PriceRequest{})
	header, err := statsSrv.Header()
	if err != nil {
		return fmt.Errorf("could not extract header: %v", err)
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
		result, err := buildBalanceResponse(s.args, name, price, req)
		if err != nil {
			return err
		}
		if err := srv.Send(result); err != nil {
			return err
		}
	}
}

func main() {
	args := parseArguments()

	// Dial account server.
	accountConn, err := grpc.Dial(args.accountServer, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v.", err)
	}
	defer accountConn.Close()
	accountClient := accountpb.NewAccountClient(accountConn)

	// Dial stats server.
	statsConn, err := grpc.Dial(args.statsServer, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v.", err)
	}
	defer statsConn.Close()
	statsClient := statspb.NewStatsClient(statsConn)

	// Listen & serve.
	lis, err := net.Listen("tcp", ":"+args.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	walletpb.RegisterWalletServer(s, &server{
		args:          &args,
		hostname:      utility.GenHostname(args.hostnameSuffix),
		accountClient: accountClient,
		statsClient:   statsClient,
	})
	reflection.Register(s)
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, healthServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
