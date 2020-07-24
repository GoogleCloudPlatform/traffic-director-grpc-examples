Disclaimer: This is not an officially supported Google product.

# gRPC examples: grpc-wallet

This repository contains example services for gRPC.

In this example, users have a wallet for a special currency: gRPC-Coin.

## wallet service

This is the user facing service, where users can fetch number of coins in
their wallets.

## stats service

This services provides price for gRPC-Coin. Users can directly query real
time coin price. The wallet service also queries coin price to calculate
total balancer in users' accounts.

## account service

This is a backend service used by wallet service and stats service, to query
user information.
