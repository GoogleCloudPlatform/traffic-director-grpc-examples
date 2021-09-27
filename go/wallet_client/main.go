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
	"os"
	"time"

	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	walletpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet"
	statspb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/stats"
	"google.golang.org/grpc/grpc-wallet/observability"
	"google.golang.org/grpc/grpc-wallet/utility"
	"google.golang.org/grpc/metadata"

	_ "google.golang.org/grpc/xds" // To enable xds support.
)

var users = map[string]map[string]string{
	"Alice": {
		"authorization": "2bd806c9",
		"membership":    "premium",
	},
	"Bob": {
		"authorization": "81b637d8",
		"membership":    "normal",
	},
}

// TODO(easwars): Get rid of this struct and define flags as package globals.
type arguments struct {
	subcommand       string
	walletServer     string
	statsServer      string
	user             string
	watch            bool
	unaryWatch       bool
	gcpClientProject string
	route            string
	creds            string
	affinity         bool
}

var args arguments

// parseArguments parses the command line arguments using the flag package.
func parseArguments() {
	flags := flag.NewFlagSet("wallet_client", flag.ExitOnError)
	flags.StringVar(&args.walletServer, "wallet_server", "localhost:18881", "address of the wallet service, default 'localhost:18881'")
	flags.StringVar(&args.statsServer, "stats_server", "localhost:18882", "address of the stats service, default 'localhost:18882'")
	flags.StringVar(&args.user, "user", "Alice", "the name of the user account, default 'Alice'")
	flags.BoolVar(&args.watch, "watch", false, "if the balance/price should be watched rather than queried once, default false")
	flags.BoolVar(&args.unaryWatch, "unary_watch", false, "if the balance/price should be watched but with repeated unary RPCs rather than a streaming rpc, default false")
	flags.StringVar(&args.gcpClientProject, "gcp_client_project", "", "if set, metrics and traces will be sent to Cloud Monitoring and Cloud Trace")
	flags.StringVar(&args.route, "route", "", "a string value to set for the 'route' header, unset by default")
	flags.StringVar(&args.creds, "creds", "insecure", "type of transport credentials to use. Supported values include 'xds' and 'insecure', defaults to 'insecure'")
	flags.BoolVar(&args.affinity, "affinity", false, "if the requests will be sent with session affinity. defaults false")
	flags.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), `
subcommands:
  balance
	print balance of the wallet
  price
	print price of grpc-coin

flags
`)
		flags.PrintDefaults()
	}

	if len(os.Args) < 2 {
		flags.Usage()
		log.Fatalln("no subcommand.")
	}
	args.subcommand = os.Args[1]
	if args.subcommand != "balance" && args.subcommand != "price" {
		flags.Usage()
		log.Fatalf("unrecognized subcommand '%s'.", args.subcommand)
	}
	flags.Parse(os.Args[2:])
	if args.user != "Alice" && args.user != "Bob" {
		flags.Usage()
		log.Fatalf("unrecognized user '%s'.", args.user)
	}
	if args.watch && args.unaryWatch {
		flags.Usage()
		log.Fatalln("unary_watch incompatible with watch.")
	}
	if args.subcommand == "price" && args.unaryWatch {
		flags.Usage()
		log.Fatalln("unary_watch incompatible with price subcommand.")
	}
}

const affinityMDKey = "session_id"

var affinityMDValues = map[string]string{
	"Alice": "11111111",
	"Bob":   "22222222",
}

// createMetaData creates the metadata for an outgoing request based on the
// flags.
func createMetaData() (metadata.MD, error) {
	md, ok := users[args.user]
	if !ok {
		return nil, fmt.Errorf("unrecognized user: %v", args.user)
	}
	if args.route != "" {
		md["route"] = args.route
	}
	if args.affinity {
		md[affinityMDKey] = affinityMDValues[args.user]
	}
	return metadata.New(md), nil
}

// handleBalanceResponse prints the data in a BalanceResponse.
func handleBalanceResponse(r *walletpb.BalanceResponse) {
	log.Printf("user: %s, total grpc-coin balance: %d.", args.user, r.GetBalance())
	for _, addrBalance := range r.GetAddresses() {
		log.Printf(" - address: %s, balance: %d.", addrBalance.Address, addrBalance.Balance)
	}
}

// balanceSubcommand handles the 'balance' subcommand.
func balanceSubcommand(opts []grpc.DialOption) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dialCancel()
	conn, err := grpc.DialContext(dialCtx, args.walletServer, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to the Wallet server: %v", err)
	}
	defer conn.Close()

	c := walletpb.NewWalletClient(conn)
	md, err := createMetaData()
	if err != nil {
		log.Fatalf("error creating metadata: %v.", err)
	}

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	if !args.watch {
		for {
			var header metadata.MD
			r, err := c.FetchBalance(ctx, &walletpb.BalanceRequest{IncludeBalancePerAddress: true}, grpc.Header(&header))
			if err != nil {
				log.Printf("failed to fetch balance: %v", err)
			} else {
				utility.PrintHostname(header)
				handleBalanceResponse(r)
			}
			if !args.unaryWatch {
				break
			}
			time.Sleep(time.Second)
		}
		return
	}
	s, err := c.WatchBalance(ctx, &walletpb.BalanceRequest{IncludeBalancePerAddress: true}, grpc.WaitForReady(true))
	if err != nil {
		log.Fatalf("failed to create stream: %v.", err)
	}
	header, err := s.Header()
	if err != nil {
		log.Fatalf("could not extract header: %v", err)
	}
	utility.PrintHostname(header)
	for {
		r, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to fetch balance: %v", err)
		}
		handleBalanceResponse(r)
	}
}

// priceSubcommand handles the 'price' subcommand.
func priceSubcommand(opts []grpc.DialOption) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dialCancel()
	conn, err := grpc.DialContext(dialCtx, args.statsServer, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to the Wallet server: %v", err)
	}
	defer conn.Close()

	c := statspb.NewStatsClient(conn)
	md, err := createMetaData()
	if err != nil {
		log.Fatalf("error creating metadata: %v.", err)
	}

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	if !args.watch {
		var header metadata.MD
		r, err := c.FetchPrice(ctx, &statspb.PriceRequest{}, grpc.Header(&header))
		if err != nil {
			log.Printf("failed to fetch price: %v", err)
		}
		utility.PrintHostname(header)
		log.Printf("grpc-coin price: %d.", r.GetPrice())
		return
	}
	s, err := c.WatchPrice(ctx, &statspb.PriceRequest{})
	if err != nil {
		log.Fatalf("failed to create stream: %v.", err)
	}
	header, err := s.Header()
	if err != nil {
		log.Fatalf("could not extract header: %v", err)
	}
	utility.PrintHostname(header)
	for {
		r, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to fetch price: %v", err)
		}
		log.Printf("grpc-coin price: %d.", r.GetPrice())
	}
}

func main() {
	parseArguments()

	// Parse credentials type from the command line to determine if xDS
	// credentials are to be used.
	xdsCreds, err := utility.ParseCredentialsType(args.creds)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("xDS credentials in use: %v", xdsCreds)

	// Use insecure credentials by default. But if xDS credentials are specified
	// on the command line, use xDS credentials with an insecure fallback.
	creds := insecure.NewCredentials()
	if xdsCreds {
		var err error
		creds, err = xdscreds.NewClientCredentials(xdscreds.ClientOptions{FallbackCreds: insecure.NewCredentials()})
		if err != nil {
			log.Fatalf("Failed to create client xDS credentials: %v", err)
		}
		log.Print("xDS credentials created")
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds), grpc.WithBlock()}

	if args.gcpClientProject != "" {
		opts = append(opts, grpc.WithStatsHandler(new(ocgrpc.ClientHandler)))
		sd := observability.ConfigureStackdriver(args.gcpClientProject)
		defer sd.Flush()
		defer sd.StopMetricsExporter()
	}

	if args.subcommand == "balance" {
		balanceSubcommand(opts)
		return
	}
	priceSubcommand(opts)
}
