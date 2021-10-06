#! /bin/bash

set -euxo pipefail

. ./00-common-env.sh

function create_private_ca_resources {
  # Create a ROOT CA.
  gcloud privateca pools create ${ROOT_CA_POOL_NAME} \
    --location ${ROOT_CA_POOL_LOCATION} \
    --tier enterprise
  gcloud privateca roots create ${ROOT_CA_NAME} \
    --pool ${ROOT_CA_POOL_NAME}
    --subject "CN=${ROOT_CA_NAME}, O=${ROOT_CA_ORGANIZATION}" \
    --max-chain-length=1 \
    --location ${ROOT_CA_POOL_LOCATION} \

  # Create a subordinate CA.
  gcloud privateca pools create ${SUBORDINATE_CA_POOL_NAME} \
    --location ${SUBORDINATE_CA_POOL_LOCATION} \
    --tier devops
  gcloud privateca subordinates create ${SUBORDINATE_CA_NAME} \
    --pool ${SUBORDINATE_CA_POOL_NAME} \
    --location ${SUBORDINATE_CA_POOL_LOCATION} \
    --issuer-pool ${ROOT_CA_POOL_NAME} \
    --issuer-location ${ROOT_CA_POOL_LOCATION} \
    --subject "CN=SUBORDINATE_CA_NAME, O=SUBORDINATE_CA_ORGANIZATION" \

  # Grant PrivateCA admin to yourself so that you can grant other privileges to
  # the GKE service account.
  gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="user:${ME}" \
    --role=roles/privateca.admin

  # Grant read permissions to the GKE service account on root CA resources.
  # This enables it to read the root CA's certificates, used to validate the
  # certificate provided by workloads.
  gcloud privateca roots add-iam-policy-binding ${ROOT_CA_NAME} \
    --location "${ROOT_CA_LOCATION}" \
    --role roles/privateca.auditor \
    --member="serviceAccount:${SA_GKE}"

  # Grant certificate manager role to the GKE service account on subordinate CA
  # resources. This enables it to submit CSRs to the subordinate CA.
  gcloud privateca roots add-iam-policy-binding "${SUBORDINATE_CA_NAME}" \
    --location "${SUBORDINATE_CA_LOCATION}" \
    --role roles/privateca.certificateManager \
    --member="serviceAccount:${SA_GKE}"

  # Apply the workload certificate configuration which sets up the trust
  # domain, the issuing CA (subordinate CA, in our case), certificate
  # properties etc.
  envsubst < WorkloadCertificateConfig.yaml | kubectl apply -f -

  # Apply the trust configuration which sets up the trust domain and the root
  # CA for validating workload certificates.
  envsubst < TrustConfig.yaml | kubectl apply -f -
}

function delete_private_ca_resources {
  kubectl delete trustconfig default
  kubectl delete workloadcertificateconfig default

  gcloud privateca roots remove-iam-policy-binding "${SUBORDINATE_CA_NAME}" \
    --location "${SUBORDINATE_CA_LOCATION}" \
    --role roles/privateca.certificateManager \
    --member="serviceAccount:${SA_GKE}"

  gcloud privateca roots remove-iam-policy-binding ${ROOT_CA_NAME} \
    --location "${ROOT_CA_LOCATION}" \
    --role roles/privateca.auditor \
    --member="serviceAccount:${SA_GKE}"

  gcloud privateca subordinates disable ${SUBORDINATE_CA_NAME} \
    --location ${SUBORDINATE_CA_LOCATION} -q
  gcloud privateca subordinates delete ${SUBORDINATE_CA_NAME} \
    --ignore-active-certificates \
    --location ${SUBORDINATE_CA_LOCATION} -q

  gcloud privateca roots disable ${ROOT_CA_NAME} \
    --location ${ROOT_CA_LOCATION} -q
  gcloud privateca roots delete ${ROOT_CA_NAME} \
    --ignore-active-certificates \
    --location ${ROOT_CA_LOCATION} -q
}
