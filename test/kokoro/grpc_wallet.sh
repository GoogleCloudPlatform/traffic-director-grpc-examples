#!/bin/bash

set -exu -o pipefail
if [[ -f /VERSION ]]; then
  cat /VERSION
fi

cd github/traffic-director-grpc-examples

pushd java
./gradlew build
popd

# Download the latest Go version in tmpdir and modify $PATH to include the
# extracted `go` binary.
walletBaseDir=${PWD}
cd ${TMPDIR}
wget https://dl.google.com/go/go1.16.5.linux-amd64.tar.gz
tar -xvf go1.16.5.linux-amd64.tar.gz
export GOROOT=${PWD}/go
export PATH="${PWD}/go/bin:${PATH}"
cd ${walletBaseDir}

pushd go
go build google.golang.org/grpc/grpc-wallet/...
popd

tools/bazel build cpp/...
