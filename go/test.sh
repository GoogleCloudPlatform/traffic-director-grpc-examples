#!/bin/bash
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

# Build
cd account_server
go build
cd ../stats_server
go build
cd ../wallet_client
go build
cd ../wallet_server
go build
cd ..

# Run servers
nohup ./account_server/account_server -hostname_suffix account > account.log 2>&1 < /dev/null &
ACCOUNT_PID="$!"
nohup ./stats_server/stats_server -hostname_suffix stats > stats.log 2>&1 < /dev/null &
STATS_PID="$!"
nohup ./wallet_server/wallet_server -hostname_suffix wallet > wallet.log 2>&1 < /dev/null &
WALLET_PID="$!"

# Run client
./wallet_client/wallet_client balance -user Alice > client.log 2>&1
./wallet_client/wallet_client balance -user Bob > client.log 2>&1
./wallet_client/wallet_client price -user Alice > client.log 2>&1
./wallet_client/wallet_client price -user Bob > client.log 2>&1

# Kill servers
kill -9 "$ACCOUNT_PID" "$STATS_PID" "$WALLET_PID"

# Echo logs for visual confirmation of output
echo "account server log:"
cat account.log
echo "stats server log:"
cat stats.log
echo "wallet server log:"
cat wallet.log
echo "wallet client log:"
cat client.log

# Clean
rm -f account.log stats.log wallet.log client.log
rm -f account_server/account_server stats_server/stats_server wallet_client/wallet_client wallet_server/wallet_server