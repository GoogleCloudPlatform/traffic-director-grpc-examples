function new_health_check() {
    gcloud compute health-checks create grpc grpcwallet-health-check \
    --use-serving-port

    gcloud compute firewall-rules create grpcwallet-allow-health-checks \
    --network default --action allow --direction INGRESS \
    --source-ranges 35.191.0.0/16,130.211.0.0/22 \
    --target-tags allow-health-checks \
    --rules tcp:50051-50053
}

function new_service() {
    # $1 = language
    # $2 = service type
    # $3 = port number
    # $4 = hostname suffix
    # $5 = additional command line arguments
    typ=$2
    port=$3
    hostname_suffix=$4
    arguments=$5

    case $1 in
        go)
            build_script="cd traffic-director-grpc-examples-master/go/${typ}_server/
sudo apt-get install -y golang git
go build ."
            server="./${typ}_server"
            ;;
        java)
            build_script="cd traffic-director-grpc-examples-master/java
sudo apt-get install -y openjdk-11-jdk-headless
./gradlew installDist"
            server="./build/install/wallet/bin/${typ}-server"
            ;;
        *)
            echo "undefined language"
            exit 123
            ;;
    esac

    startup_script="#! /bin/bash
set -ex
cd /root
export HOME=/root
sudo apt-get update -y
curl -L https://storage.googleapis.com/traffic-director/td-grpc-bootstrap-0.9.0.tar.gz | tar -xz
./td-grpc-bootstrap-0.9.0/td-grpc-bootstrap | tee /root/td-grpc-bootstrap.json
curl -L https://github.com/GoogleCloudPlatform/traffic-director-grpc-examples/archive/master.tar.gz | tar -xz
${build_script}
sudo systemd-run -E GRPC_XDS_BOOTSTRAP=/root/td-grpc-bootstrap.json ${server} --port=${port} --hostname_suffix=${hostname_suffix} ${arguments}"

    gcloud compute instance-templates create grpcwallet-${hostname_suffix}-template \
      --scopes=https://www.googleapis.com/auth/cloud-platform \
      --tags=allow-health-checks \
      --image-family=debian-10 \
      --image-project=debian-cloud \
      --metadata-from-file=startup-script=<(echo "${startup_script}")

     gcloud compute instance-groups managed create grpcwallet-${hostname_suffix}-mig-us-central1 \
       --zone us-central1-a \
       --size=2 \
       --template=grpcwallet-${hostname_suffix}-template

    gcloud compute instance-groups set-named-ports grpcwallet-${hostname_suffix}-mig-us-central1 \
      --named-ports=grpcwallet-${typ}-port:${port} \
      --zone us-central1-a 

    gcloud compute backend-services create grpcwallet-${hostname_suffix}-service \
        --global \
        --load-balancing-scheme=INTERNAL_SELF_MANAGED \
        --protocol=GRPC \
        --port-name=grpcwallet-${typ}-port \
        --health-checks grpcwallet-health-check

    gcloud compute backend-services add-backend grpcwallet-${hostname_suffix}-service \
      --instance-group grpcwallet-${hostname_suffix}-mig-us-central1 \
      --instance-group-zone us-central1-a \
      --global
}


function new_td_resources() {
    # TD resources.
    export PROJECT_ID=$(gcloud config list --format 'value(core.project)')
    gcloud compute url-maps import grpcwallet-url-map << EOF
defaultService: projects/$PROJECT_ID/global/backendServices/grpcwallet-account-service
name: grpcwallet-url-map

hostRules:
- hosts:
  - account.grpcwallet.io
  pathMatcher: grpcwallet-account-path-matcher
- hosts:
  - stats.grpcwallet.io
  pathMatcher: grpcwallet-stats-path-matcher
- hosts:
  - wallet.grpcwallet.io
  pathMatcher: grpcwallet-wallet-path-matcher

pathMatchers:
- defaultService: projects/$PROJECT_ID/global/backendServices/grpcwallet-account-service
  name: grpcwallet-account-path-matcher

- defaultService: projects/$PROJECT_ID/global/backendServices/grpcwallet-stats-service
  name: grpcwallet-stats-path-matcher
  routeRules:
  - matchRules:
    - prefixMatch: /
      headerMatches:
      - headerName: membership
        exactMatch: premium
    priority: 0
    service: projects/$PROJECT_ID/global/backendServices/grpcwallet-stats-premium-service

- defaultService: projects/$PROJECT_ID/global/backendServices/grpcwallet-wallet-v1-service
  name: grpcwallet-wallet-path-matcher
  routeRules:
  - matchRules:
    - fullPathMatch: /grpc.examples.wallet.Wallet/FetchBalance
    priority: 0
    routeAction:
      weightedBackendServices:
      - backendService: projects/$PROJECT_ID/global/backendServices/grpcwallet-wallet-v2-service
        weight: 40
      - backendService: projects/$PROJECT_ID/global/backendServices/grpcwallet-wallet-v1-service
        weight: 60
  - matchRules:
    - prefixMatch: /grpc.examples.wallet.Wallet/
    priority: 1
    routeAction:
      weightedBackendServices:
      - backendService: projects/$PROJECT_ID/global/backendServices/grpcwallet-wallet-v2-service
        weight: 100
EOF

    gcloud compute target-grpc-proxies create grpcwallet-proxy \
    --url-map grpcwallet-url-map
    #TODO:   --validate-for-proxyless

    gcloud compute forwarding-rules create grpcwallet-forwarding-rule \
    --global \
    --load-balancing-scheme=INTERNAL_SELF_MANAGED \
    --address=0.0.0.0 --address-region=us-central1 \
    --target-grpc-proxy=grpcwallet-proxy \
    --ports 80 \
    --network default
}
