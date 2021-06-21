package(default_visibility = ["//visibility:public"])

cc_library(
    name = "utility",
    srcs = [
        "utility.cc",
        "utility.h",
    ],
    defines = ["BAZEL_BUILD"],
    deps = [
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpcpp_admin",
        "@com_google_absl//absl/strings",
    ],
)

cc_binary(
    name = "client",
    srcs = ["client.cc"],
    defines = ["BAZEL_BUILD"],
    deps = [
        ":utility",
        "//:stats_proto",
        "//:wallet_proto",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc_opencensus_plugin",
        "@io_opencensus_cpp//opencensus/exporters/stats/stackdriver:stackdriver_exporter",
        "@io_opencensus_cpp//opencensus/exporters/trace/stackdriver:stackdriver_exporter",
        "@io_opencensus_cpp//opencensus/stats",
        "@io_opencensus_cpp//opencensus/trace",
    ],
)

cc_binary(
    name = "wallet-server",
    srcs = ["wallet_server.cc"],
    defines = ["BAZEL_BUILD"],
    deps = [
        ":utility",
        "//:account_proto",
        "//:stats_proto",
        "//:wallet_proto",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc++_reflection",
        "@com_github_grpc_grpc//:grpc_opencensus_plugin",
        "@io_opencensus_cpp//opencensus/exporters/stats/stackdriver:stackdriver_exporter",
        "@io_opencensus_cpp//opencensus/exporters/trace/stackdriver:stackdriver_exporter",
        "@io_opencensus_cpp//opencensus/stats",
        "@io_opencensus_cpp//opencensus/trace",
        "@io_opencensus_cpp//opencensus/trace:with_span",
    ],
)

cc_binary(
    name = "stats-server",
    srcs = ["stats_server.cc"],
    defines = ["BAZEL_BUILD"],
    deps = [
        ":utility",
        "//:account_proto",
        "//:stats_proto",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc++_reflection",
        "@com_github_grpc_grpc//:grpc_opencensus_plugin",
        "@io_opencensus_cpp//opencensus/exporters/stats/stackdriver:stackdriver_exporter",
        "@io_opencensus_cpp//opencensus/exporters/trace/stackdriver:stackdriver_exporter",
        "@io_opencensus_cpp//opencensus/stats",
        "@io_opencensus_cpp//opencensus/trace",
        "@io_opencensus_cpp//opencensus/trace:with_span",
    ],
)

cc_binary(
    name = "account-server",
    srcs = ["account_server.cc"],
    defines = ["BAZEL_BUILD"],
    deps = [
        ":utility",
        "//:account_proto",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc++_reflection",
        "@com_github_grpc_grpc//:grpc_opencensus_plugin",
        "@io_opencensus_cpp//opencensus/exporters/stats/stackdriver:stackdriver_exporter",
        "@io_opencensus_cpp//opencensus/exporters/trace/stackdriver:stackdriver_exporter",
        "@io_opencensus_cpp//opencensus/stats",
        "@io_opencensus_cpp//opencensus/trace",
        "@io_opencensus_cpp//opencensus/trace:with_span",
    ],
)
