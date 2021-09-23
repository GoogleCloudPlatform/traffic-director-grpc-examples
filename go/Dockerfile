# Dockerfile for building an image with all binaries required for the wallet
# example. To build the image, run the following command from the
# traffic-director-grpc-examples directory:
# docker build -t <TAG> -f go/Dockerfile .

FROM golang:1.16-alpine as build

# Make a traffic-director-grpc-examples directory and copy the repo into it.
WORKDIR /traffic-director-grpc-examples
COPY . .

# Build static binaries without cgo so that we can copy just the binary in the
# final image, and can get rid of the Go compiler and other dependencies.
WORKDIR /traffic-director-grpc-examples/go/account_server
RUN go build -o account-server -tags osusergo,netgo main.go

WORKDIR /traffic-director-grpc-examples/go/stats_server
RUN go build -o stats-server -tags osusergo,netgo main.go

WORKDIR /traffic-director-grpc-examples/go/wallet_server
RUN go build -o wallet-server -tags osusergo,netgo main.go

WORKDIR /traffic-director-grpc-examples/go/wallet_client
RUN go build -o wallet-client -tags osusergo,netgo main.go

# Second stage of the build which copies over only the required binaries and
# skips the Go compiler and traffic-director-grpc-examples repo from the earlier
# stage. This significantly reduces the docker image size.
FROM alpine
COPY --from=build /traffic-director-grpc-examples/go/account_server/account-server .
COPY --from=build /traffic-director-grpc-examples/go/stats_server/stats-server .
COPY --from=build /traffic-director-grpc-examples/go/wallet_server/wallet-server .
COPY --from=build /traffic-director-grpc-examples/go/wallet_client/wallet-client .
ENV GRPC_GO_LOG_VERBOSITY_LEVEL=2
ENV GRPC_GO_LOG_SEVERITY_LEVEL="info"
