#! /bin/bash

function enable_apis {
  # Enable TD and other required APIs.
  gcloud services enable \
    container.googleapis.com \
    cloudresourcemanager.googleapis.com \
    compute.googleapis.com \
    trafficdirector.googleapis.com \
    networkservices.googleapis.com \
    networksecurity.googleapis.com \
    privateca.googleapis.com

  # Allow the default GCE service account access to TD APIs.
  gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member serviceAccount:${GSA_EMAIL} \
    --role roles/trafficdirector.client
}

function create_cloud_router_instances {
  gcloud compute routers create ${CLOUD_ROUTER_NAME} \
    --network default \
    --region ${CLUSTER_REGION}

  gcloud compute routers nats create nat-config \
    --router-region ${CLUSTER_REGION} \
    --router ${CLOUD_ROUTER_NAME} \
    --nat-all-subnet-ip-ranges \
    --auto-allocate-nat-external-ips
}

function delete_cloud_router_instances {
  gcloud compute routers delete ${CLOUD_ROUTER_NAME} \
    --region ${CLUSTER_REGION} -q
}

function disable_apis {
  gcloud projects remove-iam-policy-binding ${PROJECT_ID} \
    --member serviceAccount:${GSA_EMAIL} \
    --role roles/trafficdirector.client -q
}
