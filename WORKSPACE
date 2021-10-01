workspace(name = "grpc_example")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "com_google_googleapis",
    sha256 = "150be57ff83646e5652e03683c949f0830d9a0e73ef787786864210e45537fe0",
    strip_prefix = "googleapis-6e3b55e26bf5a9f7874b6ba1411a0cc50cb87a48",
    urls = ["https://github.com/googleapis/googleapis/archive/6e3b55e26bf5a9f7874b6ba1411a0cc50cb87a48.zip"],
)

http_archive(
    name = "rules_python",
    url = "https://github.com/bazelbuild/rules_python/releases/download/0.2.0/rules_python-0.2.0.tar.gz",
    sha256 = "778197e26c5fbeb07ac2a2c5ae405b30f6cb7ad1f5510ea6fdac03bded96cc6f",
)

# Google APIs - used by Stackdriver exporter.
load("@com_google_googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
    cc = True,
    grpc = True,
)

http_archive(
    name = "com_github_grpc_grpc",
    strip_prefix = "grpc-1.41.0",
    urls = [
        "https://github.com/grpc/grpc/archive/v1.41.0.tar.gz",
    ],
)

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load(
    "@build_bazel_rules_apple//apple:repositories.bzl",
    "apple_rules_dependencies",
)

apple_rules_dependencies()
