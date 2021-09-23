# README

## Environment Variables

 * `00-common-env.sh` : This script sets common values needed by other scripts 
   as environment variables. Edit this file to set the correct values for your
   setup.

## Setup Steps

Execute the following steps in the given order:

 * Run `./init-setup.sh` to enable the required APIs and to create a NAT router.

 * Run `./gke-setup.sh` to create the server side GKE deployment (service accounts, pods etc).

 * Run `./td-setup.sh` to create the TD artifacts (backend-service and routing components).

 * Run `./security-setup.sh` to create the security policies in TD (client and server policies).

 * Finally, run `./client-setup.sh` to create the client GKE deployment. You are
   now ready to run the client as described in the next section.

## Running the Client

`ssh` into the client pod and try of the following commands:

 * The following command calls 'FetchBalance' from 'wallet-service' in a loop, 
   to demonstrate that 'FetchBalance' gets responses from 'wallet-v1' (40%)
   and 'wallet-v2' (60%).

```shell
./wallet-client balance --creds=xds --wallet_server="xds:///wallet.grpcwallet.io" --unary_watch=true
```

 * The following command calls the streaming RPC 'WatchBalance' from 'wallet-service'.
   The RPC path matches the service prefix, so all requests are sent to 'wallet-v2'.

```shell
./wallet-client balance --creds=xds --wallet_server="xds:///wallet.grpcwallet.io" --watch=true

```
 * The following commands call 'WatchPrice' from 'stats-service'. It sends the 
   user's membership (premium or not) in metadata. Premium requests are all sent
   to 'stats-premium' and get faster responses. Alice's requests always go to 
   premium and Bob's go to regular.

```shell
./wallet-client price --creds=xds --stats_server="xds:///stats.grpcwallet.io" --watch=true --user=Bob
./wallet-client price --creds=xds --stats_server="xds:///stats.grpcwallet.io" --watch=true --user=Alice
```

## Cleanup

To clean up all the artifacts created run `./cleanup.sh` 