gRPC CPP Wallet example
=======================

Build the examples:

```
../tools/bazel build :client :account-server :wallet-server :stats-server
```

Run the account server:

```
$ ../bazel-bin/cpp/account-server
```

Run the stats server:

```
$ ../bazel-bin/cpp/stats-server
```

Run the wallet server:

```
$ ../bazel-bin/cpp/wallet-server
```

Run the client:

```
$ ../bazel-bin/cpp/client balance
$ ../bazel-bin/cpp/client price
```
