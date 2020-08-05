#!/bin/bash

set -x

# Create services with different languages. Add more when more languages are supported.
./create_service.sh go   stats   50052 stats         '--account_server="xds:///account.grpcwallet.io"'
./create_service.sh go   stats   50052 stats-premium '--account_server="xds:///account.grpcwallet.io" --premium_only=true'
./create_service.sh java wallet  50051 wallet-v1     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io" --v1_behavior=true'
./create_service.sh java wallet  50051 wallet-v2     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io"'
