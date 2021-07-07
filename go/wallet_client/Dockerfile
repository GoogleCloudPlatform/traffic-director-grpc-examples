# Dockerfile for building the wallet client. To build the image, run the
# following command from the traffic-director-grpc-examples directory:
# docker build -t <TAG> -f go/wallet_client/Dockerfile .

# Use a non-alpine image as the base image so that we have some basic tools like
# bash available when we ssh into the pod to execute client commands.
FROM golang:1.16 as build

# Make a traffic-director-grpc-examples directory and copy the repo into it.
WORKDIR /traffic-director-grpc-examples
COPY . .

WORKDIR go/wallet_client
RUN go build -o wallet_client -tags osusergo,netgo main.go

## ssh into the pod and run commands

# This command calls 'FetchBalance' from 'wallet-service' in a loop,
# to demonstrate that 'FetchBalance' gets responses from 'wallet-v1' (40%)
# and 'wallet-v2' (60%).
# ./wallet_client balance --wallet_server="xds:///wallet.grpcwallet.io" --unary_watch --creds="xds"

# This command calls the streaming RPC 'WatchBalance' from 'wallet-service'.
# The RPC path matches the service prefix, so all requests
# are sent to 'wallet-v2'.
# ./wallet_client balance --wallet_server="xds:///wallet.grpcwallet.io" --watch --creds="xds"

# This command calls 'WatchPrice' from 'stats-service'. It sends the
# user's membership (premium or not) in metadata. Premium requests are
# all sent to 'stats-premium' and get faster responses. Alice's requests
# always go to premium and Bob's go to regular.
# ./wallet_client price --stats_server="xds:///stats.grpcwallet.io" --watch --user=Bob --creds="xds"
# ./wallet_client price --stats_server="xds:///stats.grpcwallet.io" --watch --user=Alice --creds="xds"
