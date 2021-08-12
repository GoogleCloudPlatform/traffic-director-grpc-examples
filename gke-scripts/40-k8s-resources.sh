#! /bin/bash

function create_k8s_resources {
  if [ $# -lt 4 ]
  then
    echo "usage: create_k8s_resources <service_name> <service_account_name> <port> <neg_name> <image_name> <optional-args-to-pass-to-binary>"
    exit 1
  fi

  export SERVICE_NAME="$1"
  export SERVICE_ACCOUNT_NAME="$2"
  export SERVICE_PORT="$3"
  export NEG_NAME="$4"
  export SERVICE_IMAGE_NAME="$5"

  shift 5

  # Convert the optional args to be passed to the binary into a comma seperated
  # string of arguments. Moving to this format makes it easier to specify as an
  # array in the `args` field of the `container` spec in the deployment.
  # TODO(easwars): See if there is a better way to do this.
  i=1
  export SERVICE_ARGS=
  while [ $# -ne 0 ]
  do
    if [ $i -ne 1 ]
    then
      SERVICE_ARGS="${SERVICE_ARGS}, "
    fi
    SERVICE_ARGS="${SERVICE_ARGS}$1"
    shift
    i=$(( $i + 1 ))
  done
  echo $SERVICE_ARGS

  envsubst < ServiceAccount.yaml | kubectl apply -f -

  gcloud iam service-accounts add-iam-policy-binding \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${PROJECT_ID}.svc.id.goog[default/${SERVICE_ACCOUNT_NAME}]" \
    ${PROJECT_NUM}-compute@developer.gserviceaccount.com

  gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member "serviceAccount:${PROJECT_ID}.svc.id.goog[default/${SERVICE_ACCOUNT_NAME}]" \
    --role roles/trafficdirector.client

  envsubst < Service.yaml | kubectl apply -f -
  envsubst < Deployment.yaml | kubectl apply -f -

  kubectl get svc ${SERVICE_NAME} -o yaml
  kubectl get deployment ${SERVICE_NAME} -o yaml
  sleep 10

  # Creation of the negs may take a while. Ensure the negs are created before
  # attempting to attach them to the TD backend services created as part of the
  # next step.
  #
  # gcloud compute network-endpoint-groups describe ${NEG_NAME} \
  #  --zone ${CLUSTER_ZONE}
  #
  # gcloud compute network-endpoint-groups list-network-endpoints ${NEG_NAME} \
  #  --zone ${CLUSTER_ZONE}
}

function delete_k8s_resources {
  if [ $# -lt 2 ]
  then
    echo "usage: delete_k8s_resources <service_name> <service_account_name>"
    exit 1
  fi
  local SERVICE_NAME="$1"
  local SERVICE_ACCOUNT_NAME="$2"

  kubectl delete svc ${SERVICE_NAME}
  kubectl delete deployment ${SERVICE_NAME}
  kubectl delete serviceaccounts ${SERVICE_ACCOUNT_NAME}
}
