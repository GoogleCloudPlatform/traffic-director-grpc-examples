#!/bin/sh
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Because of go install and go.mod, this script only runs under go/ directory.

set -eux -o pipefail

WORKDIR=$(mktemp -d)

function finish {
  rm -rf "$WORKDIR"
}
trap finish EXIT

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

mkdir -p ${WORKDIR}/out

SOURCES=$(git ls-files --exclude-standard --cached --others "../*.proto")
for src in ${SOURCES[@]}; do
  protoc --go_out=${WORKDIR}/out --go-grpc_out=${WORKDIR}/out \
    -I"."\
    -I".."\
    ${src}
done

cp -R ${WORKDIR}/out/google.golang.org/grpc/grpc-wallet/grpc .