#!/bin/bash

# This script creates a backend service for the specified service/language.
# - managed instance group
# - backend service using the instance group
#
# Run as `./create_service.sh <lang> <service> <port> <hostname_suffix> <addition_arguments>`.

set -x

language=$1
service_type=$2
port=$3
hostname_suffix=$4
shift 4 # Remaining arguments ($@) are passed to the server binary.

case $language in
    go)
        build_script="cd \"traffic-director-grpc-examples-master/go/${service_type}_server\"
sudo apt-get install -y golang git
go build ."
        server="./${service_type}_server"
        ;;
    java)
        build_script="cd traffic-director-grpc-examples-master/java
sudo apt-get install -y openjdk-11-jdk-headless
./gradlew installDist"
        server="./build/install/wallet/bin/${service_type}-server"
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
sudo systemd-run -E GRPC_XDS_BOOTSTRAP=/root/td-grpc-bootstrap.json \"${server}\" --port=${port} --hostname_suffix=${hostname_suffix} $@"

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
    --named-ports=grpcwallet-${service_type}-port:${port} \
    --zone us-central1-a 

gcloud compute backend-services create grpcwallet-${hostname_suffix}-service \
    --global \
    --load-balancing-scheme=INTERNAL_SELF_MANAGED \
    --protocol=GRPC \
    --port-name=grpcwallet-${service_type}-port \
    --health-checks grpcwallet-health-check

gcloud compute backend-services add-backend grpcwallet-${hostname_suffix}-service \
    --instance-group grpcwallet-${hostname_suffix}-mig-us-central1 \
    --instance-group-zone us-central1-a \
    --global
