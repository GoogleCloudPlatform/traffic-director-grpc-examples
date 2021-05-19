#!/bin/bash

set -exu -o pipefail
if [[ -f /VERSION ]]; then
  cat /VERSION
fi

cd github/traffic-director-grpc-examples/scripts

sudo apt-get update && sudo apt-get --only-upgrade install google-cloud-sdk-datalab google-cloud-sdk-anthos-auth google-cloud-sdk-bigtable-emulator google-cloud-sdk-pubsub-emulator google-cloud-sdk-skaffold google-cloud-sdk-firestore-emulator google-cloud-sdk-cloud-build-local google-cloud-sdk google-cloud-sdk-kubectl-oidc google-cloud-sdk-app-engine-python google-cloud-sdk-app-engine-python-extras google-cloud-sdk-spanner-emulator google-cloud-sdk-config-connector google-cloud-sdk-minikube google-cloud-sdk-app-engine-grpc google-cloud-sdk-app-engine-java kubectl google-cloud-sdk-datastore-emulator google-cloud-sdk-app-engine-go google-cloud-sdk-cbt google-cloud-sdk-local-extract google-cloud-sdk-kpt

./all.sh go
./cleanup.sh

# pushd java
# ./gradlew build
# popd

# export GOPATH="${HOME}/gopath"
# pushd go
# pushd account_server
# go build
# popd
# pushd stats_server
# go build
# popd
# pushd wallet_client
# go build
# popd
# pushd wallet_server
# go build
# popd
# popd

# tools/bazel build cpp/...
