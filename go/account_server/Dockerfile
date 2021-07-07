# Dockerfile for building the account server. To build the image, run the
# following command from the traffic-director-grpc-examples directory:
# docker build -t <TAG> -f go/account_server/Dockerfile .

FROM golang:1.16-alpine as build

# Make a traffic-director-grpc-examples directory and copy the repo into it.
WORKDIR /traffic-director-grpc-examples
COPY . .

# Build a static binary without cgo so that we can copy just the binary in the
# final image, and can get rid of the Go compiler and other dependencies.
WORKDIR go/account_server
RUN go build -o account_server -tags osusergo,netgo main.go

# Second stage of the build which copies over only the server binary and skips
# the Go compiler and traffic-director-grpc-examples repo from the earlier
# stage. This significantly reduces the docker image size.
FROM alpine
COPY --from=build /traffic-director-grpc-examples/go/account_server/account_server .
ENV GRPC_GO_LOG_VERBOSITY_LEVEL=2
ENV GRPC_GO_LOG_SEVERITY_LEVEL="info"
ENTRYPOINT ["./account_server"]
