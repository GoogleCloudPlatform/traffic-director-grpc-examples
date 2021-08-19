#! /bin/bash

set -euxo pipefail

. ./00-common-env.sh
. ./70-security-components.sh

create_security_components
