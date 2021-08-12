#! /bin/bash

function create_security_components {
  envsubst < ServerMtlsPolicy.yaml | gcloud beta network-security \
    server-tls-policies import ${SERVER_MTLS_POLICY_NAME} --location=global

  envsubst < EndpointPolicy.yaml | gcloud beta network-services \
    endpoint-policies import ${ENDPOINT_MTLS_POLICY_NAME} --location=global

  envsubst < ClientMtlsPolicy.yaml | gcloud beta network-security \
    client-tls-policies import ${CLIENT_MTLS_POLICY_NAME} --location=global

  backend_services=("${ACCOUNT_BACKEND_SERVICE_NAME}" \
    "${STATS_BACKEND_SERVICE_NAME}" \
    "${STATS_PREMIUM_BACKEND_SERVICE_NAME}" \
    "${WALLET_V1_BACKEND_SERVICE_NAME}" \
    "${WALLET_V2_BACKEND_SERVICE_NAME}")
  for bs in ${backend_services[@]}
  do
    gcloud beta compute backend-services export ${bs} --global \
      --destination=/tmp/${bs}.yaml

    envsubst < ClientSecuritySettings.yaml | cat /tmp/${bs}.yaml - \
      >/tmp/${bs}-new.yaml

    gcloud beta compute backend-services import ${bs} --global \
      --source=/tmp/${bs}-new.yaml -q
  done
}


function delete_server_security_components {
  gcloud beta network-services endpoint-policies delete ${ENDPOINT_MTLS_POLICY_NAME} --location=global -q
  gcloud beta network-security server-tls-policies delete ${SERVER_MTLS_POLICY_NAME} --location=global -q
}

function delete_client_security_components {
  gcloud beta network-security client-tls-policies delete ${CLIENT_MTLS_POLICY_NAME} --location=global -q
}
