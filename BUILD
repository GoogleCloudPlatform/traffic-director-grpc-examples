package(default_visibility = ["//visibility:public"])

load("@com_github_grpc_grpc//bazel:grpc_build_system.bzl", "grpc_proto_library")

grpc_proto_library(
    name = "wallet_proto",
    srcs = ["proto/grpc/examples/wallet/wallet.proto"],
    use_external = True,
)

grpc_proto_library(
    name = "stats_proto",
    srcs = ["proto/grpc/examples/wallet/stats/stats.proto"],
    use_external = True,
)

grpc_proto_library(
    name = "account_proto",
    srcs = ["proto/grpc/examples/wallet/account/account.proto"],
    use_external = True,
)
