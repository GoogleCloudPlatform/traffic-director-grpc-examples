#! /bin/bash

set -euxo pipefail

. ./00-common-env.sh

function create_cluster {
  gcloud container clusters create ${CLUSTER_NAME} \
    --zone=${CLUSTER_ZONE} \
    --scopes=https://www.googleapis.com/auth/cloud-platform \
    --tags=allow-health-checks \
    --release-channel rapid \
    --image-type=cos_containerd \
    --workload-pool ${WORKLOAD_POOL} \
    --enable-mesh-certificates \
    --workload-metadata=GKE_METADATA \
    --enable-ip-alias \
    --enable-autoscaling \
    --min-nodes=3 \
    --max-nodes=10

  gcloud container clusters get-credentials ${CLUSTER_NAME} --zone=${CLUSTER_ZONE}
}


function delete_cluster {
  gcloud container clusters get-credentials ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} -q

  gcloud container clusters delete ${CLUSTER_NAME} \
    --zone=${CLUSTER_ZONE} -q
}
