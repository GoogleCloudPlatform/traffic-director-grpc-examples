#!/bin/bash

# This script creates health-check for grpc, and adds a firewall rule to allow
# health check traffic to the backends.

set -x

gcloud compute health-checks create grpc grpcwallet-health-check \
    --use-serving-port

gcloud compute firewall-rules create grpcwallet-allow-health-checks \
    --network default --action allow --direction INGRESS \
    --source-ranges 35.191.0.0/16,130.211.0.0/22 \
    --target-tags allow-health-checks \
    --rules tcp:50051-50053
