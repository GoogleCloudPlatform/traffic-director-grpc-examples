#!/bin/bash
. ./utils.sh

set -x

new_health_check

new_service go account 50053 account
new_service go stats   50052 stats         '--account_server="xds:///account.grpcwallet.io"'
new_service go stats   50052 stats-premium '--account_server="xds:///account.grpcwallet.io" --premium_only=true'
new_service go wallet  50051 wallet-v1     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io" --v1_behavior=true'
new_service go wallet  50051 wallet-v2     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io"'

new_td_resources
