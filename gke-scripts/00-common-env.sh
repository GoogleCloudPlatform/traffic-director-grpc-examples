#!/bin/bash

# Basic project level configuration
export ME="YOUR_GOOGLE_EMAIL"
export PROJECT_ID="YOUR_PROJECT_ID"
export PROJECT_NUM=$(gcloud projects describe ${PROJECT_ID} --format="value(projectNumber)")
export GSA_EMAIL=${PROJECT_NUM}-compute@developer.gserviceaccount.com

# Cloud NAT configuration
export CLOUD_ROUTER_NAME="nat-router"

# Common Docker image: uncommented line is for Go. Comment the line for other languages
export WALLET_DOCKER_IMAGE="gcr.io/trafficdirector-prod/td-grpc-wallet-example-go:v2.x"

# Uncomment the line for Java
#export WALLET_DOCKER_IMAGE="gcr.io/trafficdirector-prod/td-grpc-wallet-example-java:v2.x"

# Uncomment the line for C++
#export WALLET_DOCKER_IMAGE="gcr.io/trafficdirector-prod/td-grpc-wallet-example-cpp:v2.x"

# Cluster level configuration
export CLUSTER_NAME="secure-psm-cluster"   # your cluster name here
export CLUSTER_REGION="us-west2"
export CLUSTER_ZONE="us-west2-a"
export CLUSTER_URL="https://container.googleapis.com/v1/projects/${PROJECT_ID}/locations/${CLUSTER_ZONE}/clusters/${CLUSTER_NAME}"
export WORKLOAD_POOL="${PROJECT_ID}.svc.id.goog"

# Firewall rule configuration
export FIREWALL_RULE_NAME="wallet-fw"

# Kubernetes and Traffic Director configuration for the Account service
export ACCOUNT_SERVICE_NAME="account-service"
export ACCOUNT_SERVICE_HEALTH_CHECK_NAME="account-service-hc"
export ACCOUNT_SERVICE_SA_NAME="account-service-sa"
export ACCOUNT_SERVICE_PORT="50055"
export ACCOUNT_ADMIN_PORT="50056"
export ACCOUNT_NEG_NAME="account-neg"
export ACCOUNT_SERVICE_IMAGE="${WALLET_DOCKER_IMAGE}"
export ACCOUNT_BACKEND_SERVICE_NAME="account-backend-service"
export ACCOUNT_SERVER_CMD="./account-server"

# Kubernetes and Traffic Director configuration for the Stats service
export STATS_SERVICE_NAME="stats-service"
export STATS_SERVICE_HEALTH_CHECK_NAME="stats-service-hc"
export STATS_SERVICE_SA_NAME="stats-service-sa"
export STATS_SERVICE_PORT="50053"
export STATS_ADMIN_PORT="50054"
export STATS_NEG_NAME="stats-neg"
export STATS_SERVICE_IMAGE="${WALLET_DOCKER_IMAGE}"
export STATS_BACKEND_SERVICE_NAME="stats-backend-service"
export STATS_SERVER_CMD="./stats-server"

# Kubernetes and Traffic Director configuration for the Stats Premium service
export STATS_PREMIUM_SERVICE_NAME="stats-premium-service"
export STATS_PREMIUM_SERVICE_HEALTH_CHECK_NAME="stats-premium-service-hc"
export STATS_PREMIUM_SERVICE_SA_NAME="stats-premium-service-sa"
export STATS_PREMIUM_SERVICE_PORT="50053"
export STATS_PREMIUM_ADMIN_PORT="50054"
export STATS_PREMIUM_NEG_NAME="stats-premium-neg"
export STATS_PREMIUM_SERVICE_IMAGE="${WALLET_DOCKER_IMAGE}"
export STATS_PREMIUM_BACKEND_SERVICE_NAME="stats-premium-backend-service"

# Kubernetes and Traffic Director configuration for the Wallet V1 service
export WALLET_V1_SERVICE_NAME="wallet-v1-service"
export WALLET_V1_SERVICE_HEALTH_CHECK_NAME="wallet-v1-service-hc"
export WALLET_V1_SERVICE_SA_NAME="wallet-v1-service-sa"
export WALLET_V1_SERVICE_PORT="50051"
export WALLET_V1_ADMIN_PORT="50052"
export WALLET_V1_NEG_NAME="wallet-v1-neg"
export WALLET_V1_SERVICE_IMAGE="${WALLET_DOCKER_IMAGE}"
export WALLET_V1_BACKEND_SERVICE_NAME="wallet-v1-backend-service"
export WALLET_SERVER_CMD="./wallet-server"

# Kubernetes and Traffic Director configuration for the Wallet V2 service
export WALLET_V2_SERVICE_NAME="wallet-v2-service"
export WALLET_V2_SERVICE_HEALTH_CHECK_NAME="wallet-v2-service-hc"
export WALLET_V2_SERVICE_SA_NAME="wallet-v2-service-sa"
export WALLET_V2_SERVICE_PORT="50051"
export WALLET_V2_ADMIN_PORT="50052"
export WALLET_V2_NEG_NAME="wallet-v2-neg"
export WALLET_V2_SERVICE_IMAGE="${WALLET_DOCKER_IMAGE}"
export WALLET_V2_BACKEND_SERVICE_NAME="wallet-v2-backend-service"

# Traffic Director routing configuration
export URL_MAP_NAME="wallet-url-map"
export ACCOUNT_PATH_MATCHER_NAME="account-path-matcher"
export STATS_PATH_MATCHER_NAME="stats-path-matcher"
export WALLET_PATH_MATCHER_NAME="wallet-path-matcher"
export GRPC_PROXY_NAME="wallet-target-grpc-proxy"
export FORWARDING_RULE_NAME="wallet-forwarding-rule"

# Kubernetes configuration for the Wallet client
export CLIENT_DEPLOYMENT_NAME="wallet-client"
export CLIENT_SERVICE_ACCOUNT_NAME="wallet-client-sa"
export CLIENT_IMAGE="${WALLET_DOCKER_IMAGE}"

# Private CA resources
export ROOT_CA_NAME="wallet-root-ca"
export ROOT_CA_ORGANIZATION="TestCorpLLC"
export ROOT_CA_POOL_LOCATION="us-east1"
export ROOT_CA_POOL_NAME="wallet-root-ca-pool"
export ROOT_CA_POOL_URI="//privateca.googleapis.com/projects/${PROJECT_ID}/locations/${ROOT_CA_POOL_LOCATION}/caPools/${ROOT_CA_POOL_NAME}"
export SUBORDINATE_CA_NAME="wallet-subordinate-ca"
export SUBORDINATE_CA_ORGANIZATION="TestCorpLLC"
export SUBORDINATE_CA_POOL_NAME="wallet-subordinate-ca-pool"
export SUBORDINATE_CA_POOL_LOCATION="us-east1"
export SUBORDINATE_CA_POOL_URI="//privateca.googleapis.com/projects/${PROJECT_ID}/locations/${SUBORDINATE_CA_POOL_LOCATION}/caPools/${SUBORDINATE_CA_POOL_NAME}"

# Security related configuration
export SERVER_MTLS_POLICY_NAME="server-mtls-policy"
export ENDPOINT_MTLS_POLICY_NAME="ep-mtls-psms"
export CLIENT_MTLS_POLICY_NAME="client-mtls-policy"
