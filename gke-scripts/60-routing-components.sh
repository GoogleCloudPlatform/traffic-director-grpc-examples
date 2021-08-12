#! /bin/bash

function create_routing_components {
  envsubst < ./UrlMapConfig.yaml | gcloud compute url-maps import ${URL_MAP_NAME}

  gcloud compute target-grpc-proxies create ${GRPC_PROXY_NAME} \
    --url-map ${URL_MAP_NAME} \
    --validate-for-proxyless

  gcloud compute forwarding-rules create ${FORWARDING_RULE_NAME} \
    --global \
    --load-balancing-scheme=INTERNAL_SELF_MANAGED \
    --address=0.0.0.0 \
    --target-grpc-proxy=${GRPC_PROXY_NAME} \
    --ports 80 \
    --address-region=${CLUSTER_ZONE} \
    --network default
}

function delete_routing_components {
  gcloud compute forwarding-rules delete ${FORWARDING_RULE_NAME} --global -q
  gcloud compute target-grpc-proxies delete ${GRPC_PROXY_NAME} -q
  gcloud compute url-maps delete ${URL_MAP_NAME} -q
}
