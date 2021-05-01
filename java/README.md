gRPC Java Wallet example
=======================

Build the examples:

```
$ ./gradlew installDist
```

Run the account server:

```
$ ./build/install/wallet/bin/account-server [--creds=insecure|xds]
```

Run the stats server:

```
$ ./build/install/wallet/bin/stats-server [--creds=insecure|xds]
```

Run the wallet server:

```
$ ./build/install/wallet/bin/wallet-server [--creds=insecure|xds]
```

Run the client:

```
$ ./build/install/wallet/bin/client balance [--creds=insecure|xds]
$ ./build/install/wallet/bin/client price [--creds=insecure|xds]
```
