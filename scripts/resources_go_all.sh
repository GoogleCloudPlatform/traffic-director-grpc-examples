#!/bin/bash
. ./utils.sh

set -x

new_health_check

new_service go account 50053 account
new_service go stats   50052 stats         '--account_server="xds:///account.grpcwallet.io"'
new_service go stats   50052 stats-premium '--account_server="xds:///account.grpcwallet.io" --premium_only=true'
new_service go wallet  50051 wallet-v1     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io" --v1_behavior=true'
new_service go wallet  50051 wallet-v2     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io"'

# TD resources.
export PROJECT_ID=$(gcloud config list --format 'value(core.project)')
gcloud compute url-maps import grpcwallet-url-map << EOF
defaultService: projects/$PROJECT_ID/global/backendServices/grpcwallet-account-service
name: grpcwallet-url-map

hostRules:
- hosts:
  - account.grpcwallet.io
  pathMatcher: grpcwallet-account-path-matcher
- hosts:
  - stats.grpcwallet.io
  pathMatcher: grpcwallet-stats-path-matcher
- hosts:
  - wallet.grpcwallet.io
  pathMatcher: grpcwallet-wallet-path-matcher

pathMatchers:
- defaultService: projects/$PROJECT_ID/global/backendServices/grpcwallet-account-service
  name: grpcwallet-account-path-matcher

- defaultService: projects/$PROJECT_ID/global/backendServices/grpcwallet-stats-service
  name: grpcwallet-stats-path-matcher
  routeRules:
  - matchRules:
    - prefixMatch: /
      headerMatches:
      - headerName: membership
        exactMatch: premium
    priority: 0
    service: projects/$PROJECT_ID/global/backendServices/grpcwallet-stats-premium-service

- defaultService: projects/$PROJECT_ID/global/backendServices/grpcwallet-wallet-v1-service
  name: grpcwallet-wallet-path-matcher
  routeRules:
  - matchRules:
    - fullPathMatch: /grpc.examples.wallet.Wallet/FetchBalance
    priority: 0
    routeAction:
      weightedBackendServices:
      - backendService: projects/$PROJECT_ID/global/backendServices/grpcwallet-wallet-v2-service
        weight: 40
      - backendService: projects/$PROJECT_ID/global/backendServices/grpcwallet-wallet-v1-service
        weight: 60
  - matchRules:
    - prefixMatch: /grpc.examples.wallet.Wallet/
    priority: 1
    routeAction:
      weightedBackendServices:
      - backendService: projects/$PROJECT_ID/global/backendServices/grpcwallet-wallet-v2-service
        weight: 100
EOF

gcloud compute target-grpc-proxies create grpcwallet-proxy \
  --url-map grpcwallet-url-map
#   --validate-for-proxyless

gcloud compute forwarding-rules create grpcwallet-forwarding-rule \
   --global \
   --load-balancing-scheme=INTERNAL_SELF_MANAGED \
   --address=0.0.0.0 --address-region=us-central1 \
   --target-grpc-proxy=grpcwallet-proxy \
   --ports 80 \
   --network default
