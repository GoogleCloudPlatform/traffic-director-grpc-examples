#!/bin/bash

# This script creates a backend service for the specified service/language.
# - managed instance group
# - backend service using the instance group
#
# Run as `./create_service.sh <lang> <service> <port> <hostname_suffix> <addition_arguments>`.

set -x

language="$1"
service_type="$2"
port="$3"
hostname_suffix="$4"
shift 4 # Remaining arguments ($@) are passed to the server binary.

# EXAMPLES_VERSION may be a branch name or a tag in the git repo.
EXAMPLES_VERSION=${EXAMPLES_VERSION-"master"}
# EXAMPLES_OWNER provides the ability to run the code from a forked repo.
EXAMPLES_OWNER=${EXAMPLES_OWNER-"GoogleCloudPlatform"}

size=2
if [ "${hostname_suffix}" = "wallet-v2" ]; then
    size=1
fi

case "${language}" in
    go)
        build_script="sudo apt-get install -y wget
wget https://dl.google.com/go/go1.16.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.16.5.linux-amd64.tar.gz
sudo cp /usr/local/go/bin/go /usr/bin/go
cd \"traffic-director-grpc-examples/go/${service_type}_server\"
go build ."
        server="./${service_type}_server"
        ;;
    java)
        build_script="cd traffic-director-grpc-examples/java
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
sudo apt-get install -y git
curl -L https://storage.googleapis.com/traffic-director/td-grpc-bootstrap-0.11.0.tar.gz | tar -xz
./td-grpc-bootstrap-0.11.0/td-grpc-bootstrap | tee /root/td-grpc-bootstrap.json
git clone -b ${EXAMPLES_VERSION} --single-branch --depth=1 https://github.com/${EXAMPLES_OWNER}/traffic-director-grpc-examples.git
${build_script}
sudo systemd-run -E GRPC_XDS_BOOTSTRAP=/root/td-grpc-bootstrap.json \"${server}\" --port=${port} --hostname_suffix=${hostname_suffix} $@"

gcloud compute instance-templates create grpcwallet-${hostname_suffix}-template \
    --scopes=https://www.googleapis.com/auth/cloud-platform \
    --tags=allow-health-checks \
    --image-family=debian-10 \
    --image-project=debian-cloud \
    --network-interface=no-address \
    --metadata-from-file=startup-script=<(echo "${startup_script}")

gcloud compute instance-groups managed create grpcwallet-${hostname_suffix}-mig-us-central1 \
    --zone us-central1-a \
    --size=${size} \
    --template=grpcwallet-${hostname_suffix}-template

gcloud compute instance-groups set-named-ports grpcwallet-${hostname_suffix}-mig-us-central1 \
    --named-ports=grpcwallet-${service_type}-port:${port} \
    --zone us-central1-a 

project_id="$(gcloud config list --format 'value(core.project)')"

backend_config="backends:
- balancingMode: UTILIZATION
  capacityScaler: 1.0
  group: projects/${project_id}/zones/us-central1-a/instanceGroups/grpcwallet-${hostname_suffix}-mig-us-central1
connectionDraining:
  drainingTimeoutSec: 0
healthChecks:
- projects/${project_id}/global/healthChecks/grpcwallet-health-check
loadBalancingScheme: INTERNAL_SELF_MANAGED
name: grpcwallet-${hostname_suffix}-service
portName: grpcwallet-${service_type}-port
protocol: GRPC"

if [ "${hostname_suffix}" = "stats" ]; then
    backend_config="${backend_config}
circuitBreakers:
  maxRequests: 1"
fi

gcloud compute backend-services import grpcwallet-${hostname_suffix}-service --global <<< "${backend_config}"

# The following block creates a backend service named grpcwallet-wallet-v1-affinity-service.
# However, we do not configure affinity for this BS inititally so that clients that don't
# support session-affinity feature can still work properly for all examples other than 
# the affinity example.
if [ "${hostname_suffix}" = "wallet-v1" ]; then
    backend_config="backends:
- balancingMode: UTILIZATION
  capacityScaler: 1.0
  group: projects/${project_id}/zones/us-central1-a/instanceGroups/grpcwallet-${hostname_suffix}-mig-us-central1
connectionDraining:
  drainingTimeoutSec: 0
healthChecks:
- projects/${project_id}/global/healthChecks/grpcwallet-health-check
loadBalancingScheme: INTERNAL_SELF_MANAGED
name: grpcwallet-${hostname_suffix}-affinity-service
portName: grpcwallet-${service_type}-port
protocol: GRPC
sessionAffinity: NONE
localityLbPolicy: ROUND_ROBIN"
    gcloud compute backend-services import grpcwallet-${hostname_suffix}-affinity-service --global <<< "${backend_config}"
fi
