#! /bin/bash

function create_health_check {
  if [ $# -lt 2 ]
  then
    echo "usage: create_health_check <health_check_name> <health_check_port>"
    exit 1
  fi

  gcloud compute health-checks create grpc $1 --enable-logging --port $2
}

function create_backend_service {
  if [ $# -lt 3 ]
  then
    echo "usage: create_backend_service <backend_service_name> <neg_name>"
    exit 1
  fi

  local BACKEND_SERVICE_NAME="$1"
  local HEALTH_CHECK_NAME="$2"
  local NEG_NAME="$3"

  gcloud compute backend-services create ${BACKEND_SERVICE_NAME} \
    --global \
    --load-balancing-scheme=INTERNAL_SELF_MANAGED \
    --protocol=GRPC \
    --health-checks ${HEALTH_CHECK_NAME}

  gcloud compute backend-services add-backend ${BACKEND_SERVICE_NAME} \
    --global \
    --network-endpoint-group ${NEG_NAME} \
    --network-endpoint-group-zone ${CLUSTER_ZONE} \
    --balancing-mode RATE \
    --max-rate-per-endpoint 5

  # Wait for the backend to become healthy
  gcloud compute backend-services get-health ${BACKEND_SERVICE_NAME} --global
}

function delete_backend_service {
  if [ $# -lt 2 ]
  then
    echo "usage: delete_backend_service <backend_service_name> <neg_name>"
    exit 1
  fi

  export BACKEND_SERVICE_NAME="$1"
  export NEG_NAME="$2"

  gcloud compute backend-services delete ${BACKEND_SERVICE_NAME} --global -q

  gcloud compute network-endpoint-groups delete ${NEG_NAME} \
    --zone ${CLUSTER_ZONE} -q
}

function delete_health_check {
  if [ $# -lt 1 ]
  then
    echo "usage: delete_health_check <health_check_name>"
    exit 1
  fi
  gcloud compute health-checks delete $1 -q
}
