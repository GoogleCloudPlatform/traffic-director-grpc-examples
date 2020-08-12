workspace(name = "grpc_example")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "com_github_grpc_grpc",
    urls = [
        "https://github.com/grpc/grpc/archive/4d3f0604eb5e20e94a01a1de2ac6ca27ab41e180.tar.gz",
    ],
    strip_prefix = "grpc-4d3f0604eb5e20e94a01a1de2ac6ca27ab41e180",
)

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load(
    "@build_bazel_rules_apple//apple:repositories.bzl",
    "apple_rules_dependencies",
)

apple_rules_dependencies()
