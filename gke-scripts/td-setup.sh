#! /bin/bash

set -euxo pipefail

source ./00-common-env.sh
source ./50-td-components.sh
source ./60-routing-components.sh

create_health_check ${ACCOUNT_SERVICE_HEALTH_CHECK_NAME} ${ACCOUNT_ADMIN_PORT}
create_backend_service ${ACCOUNT_BACKEND_SERVICE_NAME} ${ACCOUNT_SERVICE_HEALTH_CHECK_NAME} ${ACCOUNT_NEG_NAME}

create_health_check ${STATS_SERVICE_HEALTH_CHECK_NAME} ${STATS_ADMIN_PORT}
create_backend_service ${STATS_BACKEND_SERVICE_NAME} ${STATS_SERVICE_HEALTH_CHECK_NAME} ${STATS_NEG_NAME}

create_health_check ${STATS_PREMIUM_SERVICE_HEALTH_CHECK_NAME} ${STATS_PREMIUM_ADMIN_PORT}
create_backend_service ${STATS_PREMIUM_BACKEND_SERVICE_NAME} ${STATS_PREMIUM_SERVICE_HEALTH_CHECK_NAME} ${STATS_PREMIUM_NEG_NAME}

create_health_check ${WALLET_V1_SERVICE_HEALTH_CHECK_NAME} ${WALLET_V1_ADMIN_PORT}
create_backend_service ${WALLET_V1_BACKEND_SERVICE_NAME} ${WALLET_V1_SERVICE_HEALTH_CHECK_NAME} ${WALLET_V1_NEG_NAME}

create_health_check ${WALLET_V2_SERVICE_HEALTH_CHECK_NAME} ${WALLET_V2_ADMIN_PORT}
create_backend_service ${WALLET_V2_BACKEND_SERVICE_NAME} ${WALLET_V2_SERVICE_HEALTH_CHECK_NAME} ${WALLET_V2_NEG_NAME}

gcloud compute firewall-rules create ${FIREWALL_RULE_NAME} \
    --network default --action allow --direction INGRESS \
    --source-ranges 35.191.0.0/16,130.211.0.0/22 \
    --target-tags allow-health-checks \
    --rules tcp:${ACCOUNT_ADMIN_PORT},tcp:${STATS_ADMIN_PORT},tcp:${STATS_PREMIUM_ADMIN_PORT},tcp:${WALLET_V1_ADMIN_PORT},tcp:${WALLET_V2_ADMIN_PORT}

create_routing_components
