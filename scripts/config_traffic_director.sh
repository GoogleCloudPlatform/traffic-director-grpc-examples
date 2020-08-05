#!/bin/bash

set -x

PROJECT_ID=$(gcloud config list --format 'value(core.project)')
gcloud compute url-maps import grpcwallet-url-map --source=<(sed -e "s/\${PROJECT_ID}/${PROJECT_ID}/" url_map_template.yaml)

gcloud compute target-grpc-proxies create grpcwallet-proxy \
    --url-map grpcwallet-url-map
    #TODO:   --validate-for-proxyless

gcloud compute forwarding-rules create grpcwallet-forwarding-rule \
    --global \
    --load-balancing-scheme=INTERNAL_SELF_MANAGED \
    --address=0.0.0.0 --address-region=us-central1 \
    --target-grpc-proxy=grpcwallet-proxy \
    --ports 80 \
    --network default
