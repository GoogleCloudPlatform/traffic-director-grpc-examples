#! /bin/bash

set -euxo pipefail

. ./00-common-env.sh
. ./75-client-deployment.sh

create_client_deployment
