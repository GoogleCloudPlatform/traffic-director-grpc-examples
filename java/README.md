gRPC Java Wallet example
=======================

Build the examples:

```
$ ./gradlew installDist
```

Run the account server:

```
$ ./build/install/wallet/bin/account-server
```

Run the stats server:

```
$ ./build/install/wallet/bin/stats-server
```

Run the wallet server:

```
$ ./build/install/wallet/bin/wallet-server
```

Run the client:

```
$ ./build/install/wallet/bin/client balance
$ ./build/install/wallet/bin/client price
```
