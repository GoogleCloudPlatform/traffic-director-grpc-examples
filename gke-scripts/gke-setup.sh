#! /bin/bash

set -euxo pipefail

. ./00-common-env.sh
. ./40-k8s-resources.sh

create_k8s_resources ${ACCOUNT_SERVICE_NAME} \
  ${ACCOUNT_SERVICE_SA_NAME} \
  ${ACCOUNT_SERVICE_PORT} \
  ${ACCOUNT_NEG_NAME} \
  ${ACCOUNT_SERVICE_IMAGE} \
  ${ACCOUNT_SERVER_CMD} \
  --port=${ACCOUNT_SERVICE_PORT} \
  --admin_port=${ACCOUNT_ADMIN_PORT} \
  --creds="xds" \
  --hostname_suffix=account

create_k8s_resources ${STATS_SERVICE_NAME} \
  ${STATS_SERVICE_SA_NAME} \
  ${STATS_SERVICE_PORT} \
  ${STATS_NEG_NAME} \
  ${STATS_SERVICE_IMAGE} \
  ${STATS_SERVER_CMD} \
  --port=${STATS_SERVICE_PORT} \
  --admin_port=${STATS_ADMIN_PORT} \
  --hostname_suffix=stats \
  --creds="xds" \
  --account_server="xds:///account.grpcwallet.io"

create_k8s_resources ${STATS_PREMIUM_SERVICE_NAME} \
  ${STATS_PREMIUM_SERVICE_SA_NAME} \
  ${STATS_PREMIUM_SERVICE_PORT} \
  ${STATS_PREMIUM_NEG_NAME} \
  ${STATS_PREMIUM_SERVICE_IMAGE} \
  ${STATS_SERVER_CMD} \
  --port=${STATS_SERVICE_PORT} \
  --admin_port=${STATS_ADMIN_PORT} \
  --hostname_suffix=stats_premium \
  --creds="xds" \
  --account_server="xds:///account.grpcwallet.io" \
  --premium_only=true

create_k8s_resources ${WALLET_V1_SERVICE_NAME} \
  ${WALLET_V1_SERVICE_SA_NAME} \
  ${WALLET_V1_SERVICE_PORT} \
  ${WALLET_V1_NEG_NAME} \
  ${WALLET_V1_SERVICE_IMAGE} \
  ${WALLET_SERVER_CMD} \
  --port=${WALLET_V1_SERVICE_PORT} \
  --admin_port=${WALLET_V1_ADMIN_PORT} \
  --hostname_suffix=wallet_v1 \
  --v1_behavior=true \
  --creds="xds" \
  --account_server="xds:///account.grpcwallet.io" \
  --stats_server="xds:///stats.grpcwallet.io"

create_k8s_resources ${WALLET_V2_SERVICE_NAME} \
  ${WALLET_V2_SERVICE_SA_NAME} \
  ${WALLET_V2_SERVICE_PORT} \
  ${WALLET_V2_NEG_NAME} \
  ${WALLET_V2_SERVICE_IMAGE} \
  ${WALLET_SERVER_CMD} \
  --port=${WALLET_V2_SERVICE_PORT} \
  --admin_port=${WALLET_V2_ADMIN_PORT} \
  --hostname_suffix=wallet_v2 \
  --creds="xds" \
  --account_server="xds:///account.grpcwallet.io" \
  --stats_server="xds:///stats.grpcwallet.io"

