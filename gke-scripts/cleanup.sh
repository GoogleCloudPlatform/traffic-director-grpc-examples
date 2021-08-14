#! /bin/bash

set -x

. ./00-common-env.sh
. ./10-apis.sh
. ./40-k8s-resources.sh
. ./50-td-components.sh
. ./60-routing-components.sh
. ./70-security-components.sh
. ./75-client-deployment.sh

delete_client_deployment

delete_routing_components

delete_backend_service ${ACCOUNT_BACKEND_SERVICE_NAME} ${ACCOUNT_NEG_NAME}
delete_health_check ${ACCOUNT_SERVICE_HEALTH_CHECK_NAME}

delete_backend_service ${STATS_BACKEND_SERVICE_NAME} ${STATS_NEG_NAME}
delete_health_check ${STATS_SERVICE_HEALTH_CHECK_NAME}

delete_backend_service ${STATS_PREMIUM_BACKEND_SERVICE_NAME} ${STATS_PREMIUM_NEG_NAME}
delete_health_check ${STATS_PREMIUM_SERVICE_HEALTH_CHECK_NAME}

delete_backend_service ${WALLET_V1_BACKEND_SERVICE_NAME} ${WALLET_V1_NEG_NAME}
delete_health_check ${WALLET_V1_SERVICE_HEALTH_CHECK_NAME}

delete_backend_service ${WALLET_V2_BACKEND_SERVICE_NAME} ${WALLET_V2_NEG_NAME}
delete_health_check ${WALLET_V2_SERVICE_HEALTH_CHECK_NAME}

delete_k8s_resources ${ACCOUNT_SERVICE_NAME} ${ACCOUNT_SERVICE_SA_NAME}
delete_k8s_resources ${STATS_SERVICE_NAME} ${STATS_SERVICE_SA_NAME}
delete_k8s_resources ${STATS_PREMIUM_SERVICE_NAME} ${STATS_PREMIUM_SERVICE_SA_NAME}
delete_k8s_resources ${WALLET_V1_SERVICE_NAME} ${WALLET_V1_SERVICE_SA_NAME}
delete_k8s_resources ${WALLET_V2_SERVICE_NAME} ${WALLET_V2_SERVICE_SA_NAME}

delete_cloud_router_instances

delete_server_security_components
delete_client_security_components

gcloud compute firewall-rules delete ${FIREWALL_RULE_NAME}  -q
delete_policy_bindings
