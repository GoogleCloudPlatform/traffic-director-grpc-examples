#!/bin/bash

set -exu -o pipefail
if [[ -f /VERSION ]]; then
  cat /VERSION
fi

cd github/traffic-director-grpc-examples

pushd java
./gradlew build
popd

export GOPATH="${HOME}/gopath"
pushd go
pushd account_server
go build
popd
pushd stats_server
go build
popd
pushd wallet_client
go build
popd
pushd wallet_server
go build
popd
popd

tools/bazel build cpp/...
