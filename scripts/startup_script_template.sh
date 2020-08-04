#! /bin/bash

# $build  = commands to build the server
# $server = path to the server binary

set -ex
cd /root
export HOME=/root
sudo apt-get update -y
curl -L https://storage.googleapis.com/traffic-director/td-grpc-bootstrap-0.9.0.tar.gz | tar -xz
./td-grpc-bootstrap-0.9.0/td-grpc-bootstrap | tee /root/td-grpc-bootstrap.json
curl -L https://github.com/GoogleCloudPlatform/traffic-director-grpc-examples/archive/master.tar.gz | tar -xz
${build}
sudo systemd-run -E GRPC_XDS_BOOTSTRAP=/root/td-grpc-bootstrap.json ${server}