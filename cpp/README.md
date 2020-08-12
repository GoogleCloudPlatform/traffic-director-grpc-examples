gRPC CPP Wallet example
=======================

Build the examples:

```
../tools/bazel build :client :account-server :wallet-server :stats-server
```

Run the account server:

```
$ ../bazel-bin/account-server
```

Run the stats server:

```
$ ../bazel-bin/stats-server
```

Run the wallet server:

```
$ ../bazel-bin/wallet-server
```

Run the client:

```
$ ../bazel-bin/client balance
$ ../bazel-bin/client price
```
