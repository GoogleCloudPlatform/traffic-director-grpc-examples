#! /bin/bash

function create_client_deployment {
  envsubst < ClientDeployment.yaml | kubectl apply -f -

  gcloud iam service-accounts add-iam-policy-binding \
    --role=roles/iam.workloadIdentityUser \
    --member="serviceAccount:${PROJECT_ID}.svc.id.goog[default/${CLIENT_SERVICE_ACCOUNT_NAME}]" \
    ${PROJECT_NUM}-compute@developer.gserviceaccount.com

  gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${PROJECT_ID}.svc.id.goog[default/${CLIENT_SERVICE_ACCOUNT_NAME}]" \
    --role=roles/trafficdirector.client

  kubectl get deployment ${CLIENT_DEPLOYMENT_NAME} -o yaml
}

function delete_client_deployment {
  gcloud iam service-accounts remove-iam-policy-binding \
    --role=roles/iam.workloadIdentityUser \
    --member="serviceAccount:${PROJECT_ID}.svc.id.goog[default/${CLIENT_SERVICE_ACCOUNT_NAME}]" \
    ${PROJECT_NUM}-compute@developer.gserviceaccount.com

  gcloud projects remove-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${PROJECT_ID}.svc.id.goog[default/${CLIENT_SERVICE_ACCOUNT_NAME}]" \
    --role=roles/trafficdirector.client

  kubectl delete serviceaccount ${CLIENT_SERVICE_ACCOUNT_NAME}
  kubectl delete deployment ${CLIENT_DEPLOYMENT_NAME}
}
