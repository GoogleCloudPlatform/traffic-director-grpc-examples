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

package utility

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/codes"
	accountpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet/account"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// PrintHostname pulls the hostname from a response header and prints it.
func PrintHostname(header metadata.MD) {
	hostnames, ok := header["hostname"]
	if !ok {
		log.Printf("server host: error: no hostname")
		return
	}
	if len(hostnames) < 1 {
		log.Printf("server host: error: no hostname")
		return
	}
	log.Printf("server host: %s", hostnames[0])
}

// GenHostname generates the hostname for a service given the hostname_suffix
// flag.
func GenHostname(hostnameSuffix string) string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "generated-" + strconv.Itoa(rand.Intn(1000))
	}
	if hostnameSuffix != "" {
		return hostname + "_" + hostnameSuffix
	}
	return hostname
}

// ValidateMembership performs an rpc to the account service to validate a
// user's membership.
func ValidateMembership(inCtx context.Context, accountClient accountpb.AccountClient) (premium bool, name, token, requestedMembership string, err error) {
	// Pull metadata from incoming context.
	inMd, ok := metadata.FromIncomingContext(inCtx)
	if !ok {
		return false, "", "", "", status.Error(codes.Unauthenticated, "missing user authentication")
	}

	// Pull token and membership from incoming metadata.
	tokens, ok := inMd["authorization"]
	if !ok || len(tokens) < 1 {
		return false, "", "", "", status.Error(codes.Unauthenticated, "missing user authentication")
	}
	memberships, ok := inMd["membership"]
	if !ok || len(memberships) < 1 {
		return false, "", "", "", status.Error(codes.Unauthenticated, "missing user authentication")
	}
	token = tokens[0]
	requestedMembership = memberships[0]
	if requestedMembership != "normal" && requestedMembership != "premium" {
		return false, "", "", "", status.Error(codes.Unauthenticated, "unknown user membership")
	}
	routes, ok := inMd["route"]
	if ok {
	  inCtx = metadata.NewOutgoingContext(inCtx, metadata.Pairs("route", routes[0]))
	}

	// Perform RPC to account client.
	var header metadata.MD
	r, err := accountClient.GetUserInfo(inCtx, &accountpb.GetUserInfoRequest{Token: token}, grpc.Header(&header))
	if err != nil {
		return false, "", "", "", status.Errorf(status.Code(err), "could not get balance: %v", err)
	}
	PrintHostname(header)

	// Parse membership from request response.
	actualMembership := ""
	switch r.GetMembership() {
	case accountpb.MembershipType_PREMIUM:
		actualMembership = "premium"
	case accountpb.MembershipType_NORMAL:
		actualMembership = "normal"
	default:
		return false, "", "", "", status.Error(codes.Unauthenticated, "unrecognized user membership type")
	}

	// Print results & return.
	name = r.GetName()
	success := requestedMembership <= actualMembership
	log.Printf("token: %s, name: %s, membership: %s, requested membership: %s, authentication success: %v", token, name, actualMembership, requestedMembership, success)
	if !success {
		return false, "", "", "", status.Error(codes.Unauthenticated, "requested membership higher than actual membership")
	}
	return requestedMembership != "normal", name, token, requestedMembership, nil
}

// ParseCredentialsType parses the value set for the command line flag -creds.
// Supported values for this flag include 'insecure' and 'xds'. When a nil error
// is returned, the first return value indicates whether xDS credentials are to
// be used. A non-nil error is returned when credsType is set to an unsupported
// credentials type.
func ParseCredentialsType(credsType string) (bool, error) {
	switch credsType {
	case "insecure":
		return false, nil
	case "xds":
		return true, nil
	default:
		return false, fmt.Errorf("-creds set to unsupported value %q. Only supported values are 'insecure' and 'xds'", credsType)
	}
}

// StartAdminServices exposes admin services (csds, health and reflection) on
// the on admin_port. It is a blocking call and will return only when the gRPC
// server exposing the admin services is stopped.
func StartAdminServices(adminPort string) error {
	adminListener, err := net.Listen("tcp", ":"+adminPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	cleanup, err := admin.Register(s)
	if err != nil {
		return fmt.Errorf("failed to register admin service: %v", err)
	}
	defer cleanup()
	reflection.Register(s)
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, healthServer)
	if err := s.Serve(adminListener); err != nil {
		return fmt.Errorf("failed to serve admin services: %v", err)
	}
	return nil
}
