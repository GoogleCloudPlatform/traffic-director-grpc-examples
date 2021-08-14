#! /bin/bash

set -euxo pipefail

. ./00-common-env.sh
. ./10-apis.sh

enable_apis
create_cloud_router_instances

