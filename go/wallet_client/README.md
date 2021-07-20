## Steps to run the wallet client on k8s

After bringing up all the services and the wallet client in k8s pods, SSH into
the client pod and run the following commands.

1. This command calls `FetchBalance` from `wallet-service` in a loop,
   to demonstrate that `FetchBalance` gets responses from `wallet-v1` (40%)
   and `wallet-v2` (60%).

   > ./wallet_client balance --wallet_server="xds:///wallet.grpcwallet.io" --unary_watch --creds="xds"

1. This command calls the streaming RPC `WatchBalance` from `wallet-service`.
   The RPC path matches the service prefix, so all requests are sent to
    `wallet-v2`.

    > ./wallet_client balance --wallet_server="xds:///wallet.grpcwallet.io" --watch --creds="xds"

 1. This command calls `WatchPrice` from `stats-service`. It sends the user's
    membership (premium or not) in metadata. Premium requests are all sent to
    `stats-premium` and get faster responses. Alice's requests always go to
     premium and Bob's go to regular.

     > ./wallet_client price --stats_server="xds:///stats.grpcwallet.io" --watch --user=Bob --creds="xds"

     > ./wallet_client price --stats_server="xds:///stats.grpcwallet.io" --watch --user=Alice --creds="xds"
