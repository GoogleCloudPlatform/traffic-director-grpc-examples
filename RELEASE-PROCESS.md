# Release Process

External facing [Traffic Director](https://cloud.google.com/traffic-director)
documentation uses the wallet example to showcase various xDS features. Here is
one such [example][wallet-example-in-td-docs] where a specific branch of the
wallet example is used.

[wallet-example-in-td-docs]:
https://cloud.google.com/traffic-director/docs/proxyless-configure-advanced-traffic-management#verify-routing-configuration

# When to make a new version
Whenever we have a new iteration of the public documentation, referencing the
wallet example to showcase new set of functionality, we must make a new branch.
Public documentation for the new set of functionality would end up referencing
the newly created branch. This would ensure that existing documentation wouldn't
need any changes, and any changes we make to the examples to support new
functionality would not end up breaking old documentation inadvertently.

# How to make a new version
Once all changes are made to examples to support the new functionality, we
should create a new branch. The public documentation will always refer to `HEAD`
of these branches. This approach has the advantage that any bug fixes that need
to be made on a branch do not require any update to the documentation.

The branch corresponding to the first release version on this repo was v1.0.x.
But going forward, we would just do v{VERSION}.x.

```
# Make a local repo.
mkdir /tmp/wallet && cd /tmp/wallet
git clone https://github.com/GoogleCloudPlatform/traffic-director-grpc-examples.git
cd traffic-director-grpc-examples

# Make the new branch.
NEXT_RELEASE_VERSION=v{{CURRENT_VERSION+1}}.x
git branch ${NEXT_RELEASE_VERSION}
git push origin ${NEXT_RELEASE_VERSION}
```
