/*
 *
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

#include <grpcpp/ext/proto_server_reflection_plugin.h>
#include <grpcpp/grpcpp.h>
#include <grpcpp/health_check_service_interface.h>
#include <grpcpp/opencensus.h>
#include <unistd.h>

#include <iostream>
#include <memory>
#include <string>

#include "opencensus/exporters/stats/stackdriver/stackdriver_exporter.h"
#include "opencensus/exporters/trace/stackdriver/stackdriver_exporter.h"
#include "opencensus/trace/with_span.h"
#include "proto/grpc/examples/wallet/account/account.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using grpc::StatusCode;
using grpc::examples::wallet::account::Account;
using grpc::examples::wallet::account::GetUserInfoRequest;
using grpc::examples::wallet::account::GetUserInfoResponse;
using grpc::examples::wallet::account::MembershipType;

class AccountServiceImpl final : public Account::Service {
 public:
  explicit AccountServiceImpl(const std::string& hostname)
      : hostname_(hostname) {}

 private:
  Status GetUserInfo(ServerContext* context, const GetUserInfoRequest* request,
                     GetUserInfoResponse* response) override {
    opencensus::trace::Span span = grpc::GetSpanFromServerContext(context);
    {
      // Run in OpenCensus span received from the client to correlate the traces
      // in Cloud Monitoring.
      opencensus::trace::WithSpan ws(span);
      std::string token = request->token();
      context->AddInitialMetadata("hostname", hostname_);
      if (token == "2bd806c9") {
        response->set_name("Alice");
        response->set_membership(MembershipType::PREMIUM);
      } else if (token == "81b637d8") {
        response->set_name("Bob");
        response->set_membership(MembershipType::NORMAL);
      } else {
        return Status(StatusCode::NOT_FOUND, "user not found");
      }
      return Status::OK;
    }
  }

  std::string hostname_;
};

void RunServer(const std::string& port, const std::string& hostname_suffix) {
  std::string hostname;
  char base_hostname[256];
  if (gethostname(base_hostname, 256) != 0) {
    sprintf(base_hostname, "%s-%d", "generated", rand() % 1000);
  }
  hostname = std::string(base_hostname);
  hostname += hostname_suffix;
  std::string server_address("0.0.0.0:");
  server_address += port;
  AccountServiceImpl service(hostname);
  grpc::EnableDefaultHealthCheckService(true);
  grpc::reflection::InitProtoReflectionServerBuilderPlugin();
  // Listen on the given address without any authentication mechanism.
  std::cout << "Account Server listening on " << server_address << std::endl;
  ServerBuilder builder;
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  // Register "service" as the instance through which we'll communicate with
  // clients. In this case it corresponds to an *synchronous* service.
  builder.RegisterService(&service);
  // Finally assemble the server.
  std::unique_ptr<Server> server(builder.BuildAndStart());
  // Wait for the server to shutdown. Note that some other thread must be
  // responsible for shutting down the server for this call to ever return.
  server->Wait();
}

int main(int argc, char** argv) {
  std::string port = "18883";
  std::string hostname_suffix = "";
  std::string observability_project = "";
  std::string arg_str_port("--port");
  std::string arg_str_hostname_suffix("--hostname_suffix");
  std::string arg_str_observability_project("--observability_project");
  for (int i = 1; i < argc; ++i) {
    std::string arg_val = argv[i];
    size_t start_pos = arg_val.find(arg_str_port);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_port.size();
      if (arg_val[start_pos] == '=') {
        port = arg_val.substr(start_pos + 1);
        continue;
      } else {
        std::cout << "The only correct argument syntax is --port=" << std::endl;
        return 1;
      }
    }
    start_pos = arg_val.find(arg_str_hostname_suffix);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_hostname_suffix.size();
      if (arg_val[start_pos] == '=') {
        hostname_suffix = arg_val.substr(start_pos + 1);
        continue;
      } else {
        std::cout << "The only correct argument syntax is --hostname_suffix="
                  << std::endl;
        return 1;
      }
    }
    start_pos = arg_val.find(arg_str_observability_project);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_observability_project.size();
      if (arg_val[start_pos] == '=') {
        observability_project = arg_val.substr(start_pos + 1);
        continue;
      } else {
        std::cout
            << "The only correct argument syntax is --observability_project="
            << std::endl;
        return 1;
      }
    }
  }
  std::cout << "Account Server arguments: port: " << port
            << ", hostname_suffix: " << hostname_suffix
            << ", observability_project: " << observability_project
            << std::endl;
  if (!observability_project.empty()) {
    grpc::RegisterOpenCensusPlugin();
    grpc::RegisterOpenCensusViewsForExport();
    opencensus::trace::TraceConfig::SetCurrentTraceParams(
        {128, 128, 128, 128, opencensus::trace::ProbabilitySampler(1.0)});
    opencensus::exporters::trace::StackdriverOptions trace_opts;
    trace_opts.project_id = observability_project;
    opencensus::exporters::trace::StackdriverExporter::Register(
        std::move(trace_opts));
    opencensus::exporters::stats::StackdriverOptions stats_opts;
    stats_opts.project_id = observability_project;
    opencensus::exporters::stats::StackdriverExporter::Register(
        std::move(stats_opts));
  }
  RunServer(port, hostname_suffix);
  return 0;
}
