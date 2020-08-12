package(default_visibility = ["//visibility:public"])

cc_binary(
    name = "client",
    srcs = ["client.cc"],
    defines = ["BAZEL_BUILD"],
    deps = [
        "//:stats_proto",
        "//:wallet_proto",
        "@com_github_grpc_grpc//:grpc++",

    ],
)

cc_binary(
    name = "wallet-server",
    srcs = ["wallet_server.cc"],
    defines = ["BAZEL_BUILD"],
    deps = [
        "//:account_proto",
        "//:stats_proto",
        "//:wallet_proto",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc++_reflection",
    ],
)

cc_binary(
    name = "stats-server",
    srcs = ["stats_server.cc"],
    defines = ["BAZEL_BUILD"],
    deps = [
        "//:account_proto",
        "//:stats_proto",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc++_reflection",
    ],
)

cc_binary(
    name = "account-server",
    srcs = ["account_server.cc"],
    defines = ["BAZEL_BUILD"],
    deps = [
        "//:account_proto",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc++_reflection",
    ],
)
