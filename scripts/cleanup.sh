#!/bin/bash
set -x

# This scripts deletes all resources created by this example.

# TD resources
gcloud compute forwarding-rules delete grpcwallet-forwarding-rule --global -q
gcloud compute target-grpc-proxies delete grpcwallet-proxy -q
gcloud compute url-maps delete grpcwallet-url-map -q

# per service
services=(
    "grpcwallet-account"
    "grpcwallet-stats"
    "grpcwallet-stats-premium"
    "grpcwallet-wallet-v1"
    "grpcwallet-wallet-v2"
)
gcloud compute backend-services delete grpcwallet-wallet-v1-affinity-service --global -q
for s in "${services[@]}"; do
    gcloud compute backend-services delete "$s"-service --global -q
    gcloud compute instance-groups managed delete "$s"-mig-us-central1 --zone us-central1-a -q
    gcloud compute instance-templates delete "$s"-template -q
done

# health check & firewall
gcloud compute firewall-rules delete grpcwallet-allow-health-checks -q
gcloud compute health-checks delete grpcwallet-health-check -q
