# Release Process

External facing [Traffic Director](https://cloud.google.com/traffic-director)
documentation for specific features tend to use the wallet example to showcase
the feature better. Here is one such [example][wallet-example-in-td-docs] where
a specific branch of the wallet example is used.

[wallet-example-in-td-docs]:
https://cloud.google.com/traffic-director/docs/proxyless-configure-advanced-traffic-management#verify-routing-configuration

# When to make a new release
Whenever we have a new iteration of the public documentation, referencing the
wallet example to showcase new set of functionality, we should consider making a
new release. Public documentation for the new set of functionality would end up
referencing the newly created version. This would ensure that existing
documentation wouldn't need any changes, and any changes we make to the examples
to support new functionality would not end up breaking old documentation
inadvertently.

# How to make a new release
Once all changes are made to examples to support the new functionality, we
should create a new branch to represent the new release.

This approach has the advantage that any bug fixes that need to be made on a
release branch do not require a patch release. The public documentation will
always refer to `HEAD` of these branches, and will not require any change to
pick up the bug fixes.

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

Use the GitHub [release
page](https://github.com/GoogleCloudPlatform/traffic-director-grpc-examples/releases/new)
to create the new release.
1. Change branch ("target") to v{{NEXT_RELEASE_VERSION}}.x
1. Tag name: "v{{NEXT_RELEASE_VERSION}}.0"
1. Title: "Release v{{NEXT_RELEASE_VERSION}}.0"
1. Body: Copy commit messages
1. Double-check release branch ("target")
1. Save as a draft release and have someone review and approve
1. Publish Release
