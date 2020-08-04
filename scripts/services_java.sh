#!/bin/bash

function new_service() {
    # $1 = service type
    # $2 = port number
    # $3 = hostname suffix
    # $4 = additional command line arguments
    typ=$1
    port=$2
    hostname_suffix=$3
    arguments=$4

    build_script="cd traffic-director-grpc-examples-master/java
sudo apt-get install -y openjdk-11-jdk-headless
./gradlew installDist"
    server="./build/install/wallet/bin/${typ}-server"

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

set -x

# This script is to be used with the examples doc [url]. The account service is
# created manually. This script only creates the remaining services.

new_service stats   50052 stats         '--account_server="xds:///account.grpcwallet.io"'
new_service stats   50052 stats-premium '--account_server="xds:///account.grpcwallet.io" --premium_only=true'
new_service wallet  50051 wallet-v1     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io" --v1_behavior=true'
new_service wallet  50051 wallet-v2     '--account_server="xds:///account.grpcwallet.io" --stats_server="xds:///stats.grpcwallet.io"'
