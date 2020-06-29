#!/bin/bash
set -e
# Clean
rm -f account.log stats.log wallet.log client.log
# Proto
cd grpc/examples/wallet
protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative wallet.proto
protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative account/account.proto
protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative stats/stats.proto
cd ../../..
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