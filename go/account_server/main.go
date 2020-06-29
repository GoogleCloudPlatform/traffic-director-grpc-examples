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

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	accountpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/account"
	"google.golang.org/grpc/grpc-wallet/utility"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type user struct {
	name       string
	membership accountpb.MembershipType
}

var users = map[string]user{
	"2bd806c9": {"Alice", accountpb.MembershipType_PREMIUM},
	"81b637d8": {"Bob", accountpb.MembershipType_NORMAL},
}

type arguments struct {
	port           string
	hostnameSuffix string
}

// parseArguments parses the command line arguments using the flag package.
func parseArguments() arguments {
	result := arguments{}
	flag.StringVar(&result.port, "port", "18883", "the port to listen on, default '18883'")
	flag.StringVar(&result.hostnameSuffix, "hostname_suffix", "", "suffix to append to hostname in response header for outgoing RPCs, default ''")
	flag.Parse()
	return result
}

type server struct {
	accountpb.UnimplementedAccountServer
	hostname string
}

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

func main() {
	args := parseArguments()

	lis, err := net.Listen("tcp", ":"+args.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	accountpb.RegisterAccountServer(s, &server{hostname: utility.GenHostname(args.hostnameSuffix)})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
