#!/bin/bash

set -x

# $1 = language, one of "java" or "go"
if ! [[ $1 =~ ^(go|java)$ ]]; then
    echo "language $1 is undefined, pick one from [go, java]"
    exit 123
fi

./create_health_check.sh
./create_service.sh $1 account 50053 account
./create_service.sh $1 stats   50052 stats         '--account_server="xds:///account.grpcwallet.io"'
./create_service.sh $1 stats   50052 stats-premium '--account_server="xds:///account.grpcwallet.io" --premium_only=true'
./create_service.sh $1 wallet  50051 wallet-v1     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io" --v1_behavior=true'
./create_service.sh $1 wallet  50051 wallet-v2     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io"'
./config_traffic_director.sh
