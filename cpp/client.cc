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

#include <unistd.h>

#include <iostream>
#include <memory>
#include <string>
#include <thread>

#include "grpc++/grpc++.h"
#include "grpcpp/opencensus.h"
#include "opencensus/exporters/stats/stackdriver/stackdriver_exporter.h"
#include "opencensus/exporters/trace/stackdriver/stackdriver_exporter.h"
#include "proto/grpc/examples/wallet/stats/stats.grpc.pb.h"
#include "proto/grpc/examples/wallet/wallet.grpc.pb.h"
#include "utility.h"

using grpc::Channel;
using grpc::ChannelArguments;
using grpc::ClientContext;
using grpc::ClientReader;
using grpc::Status;
using grpc::examples::wallet::BalanceRequest;
using grpc::examples::wallet::BalanceResponse;
using grpc::examples::wallet::Wallet;
using grpc::examples::wallet::stats::PriceRequest;
using grpc::examples::wallet::stats::PriceResponse;
using grpc::examples::wallet::stats::Stats;

class WalletClient {
 public:
  WalletClient(std::shared_ptr<Channel> channel)
      : stub_(Wallet::NewStub(channel)) {}

  void FetchBalance(const std::string& user, const std::string& route,
                    const bool affinity) {
    BalanceRequest request;
    request.set_include_balance_per_address(true);
    BalanceResponse response;
    ClientContext context;
    context.set_wait_for_ready(true);
    if (user == "Alice") {
      context.AddMetadata("authorization", "2bd806c9");
      context.AddMetadata("membership", "premium");
      if (affinity) {
        // use something unique per user as session id.
        context.AddMetadata("session_id", "11111111");
      }
    } else {
      context.AddMetadata("authorization", "81b637d8");
      context.AddMetadata("membership", "normal");
      if (affinity) {
        // use something unique per user as session id.
        context.AddMetadata("session_id", "22222222");
      }
    }
    if (route != "") {
      context.AddMetadata("route", route);
    }
    Status status = stub_->FetchBalance(&context, request, &response);
    if (status.ok()) {
      auto metadata_hostname =
          context.GetServerInitialMetadata().find("hostname");
      if (metadata_hostname != context.GetServerInitialMetadata().end()) {
        std::cout << "server host: "
                  << std::string(metadata_hostname->second.data(),
                                 metadata_hostname->second.length())
                  << std::endl;
      }
      std::cout << "user: " << user
                << " total grpc-coin balance: " << response.balance()
                << std::endl;
      for (const auto& address : response.addresses()) {
        std::cout << " - address: " << address.address()
                  << ", balance: " << address.balance() << std::endl;
      }
    } else {
      std::cout << status.error_code() << ": " << status.error_message()
                << std::endl;
    }
  }

  void WatchBalance(const std::string& user, const std::string& route) {
    BalanceRequest request;
    request.set_include_balance_per_address(true);
    BalanceResponse response;
    ClientContext context;
    context.set_wait_for_ready(true);
    if (user == "Alice") {
      context.AddMetadata("authorization", "2bd806c9");
      context.AddMetadata("membership", "premium");
    } else {
      context.AddMetadata("authorization", "81b637d8");
      context.AddMetadata("membership", "normal");
    }
    if (route != "") {
      context.AddMetadata("route", route);
    }
    std::unique_ptr<ClientReader<BalanceResponse> > reader(
        stub_->WatchBalance(&context, request));
    bool first_read = true;
    while (reader->Read(&response)) {
      if (first_read) {
        auto metadata_hostname =
            context.GetServerInitialMetadata().find("hostname");
        if (metadata_hostname != context.GetServerInitialMetadata().end()) {
          std::cout << "server host: "
                    << std::string(metadata_hostname->second.data(),
                                   metadata_hostname->second.length())
                    << std::endl;
        }
        first_read = false;
      }
      std::cout << "user: " << user
                << " total grpc-coin balance: " << response.balance()
                << std::endl;
      for (const auto& address : response.addresses()) {
        std::cout << " - address: " << address.address()
                  << ", balance: " << address.balance() << std::endl;
      }
    }
    Status status = reader->Finish();
    if (!status.ok()) {
      std::cout << status.error_code() << ": " << status.error_message()
                << std::endl;
    }
  }

 private:
  std::unique_ptr<Wallet::Stub> stub_;
};

class StatsClient {
 public:
  StatsClient(std::shared_ptr<Channel> channel)
      : stub_(Stats::NewStub(channel)) {}

  void FetchPrice(const std::string& user, const std::string& route) {
    PriceRequest request;
    PriceResponse response;
    ClientContext context;
    context.set_wait_for_ready(true);
    if (user == "Alice") {
      context.AddMetadata("authorization", "2bd806c9");
      context.AddMetadata("membership", "premium");
    } else {
      context.AddMetadata("authorization", "81b637d8");
      context.AddMetadata("membership", "normal");
    }
    if (route != "") {
      context.AddMetadata("route", route);
    }
    Status status = stub_->FetchPrice(&context, request, &response);
    if (status.ok()) {
      auto metadata_hostname =
          context.GetServerInitialMetadata().find("hostname");
      if (metadata_hostname != context.GetServerInitialMetadata().end()) {
        std::cout << "server host: "
                  << std::string(metadata_hostname->second.data(),
                                 metadata_hostname->second.length())
                  << std::endl;
      }
      std::cout << "grpc-coin price: " << response.price() << std::endl;
    } else {
      std::cout << status.error_code() << ": " << status.error_message()
                << std::endl;
    }
  }

  void WatchPrice(const std::string& user, const std::string& route) {
    PriceRequest request;
    PriceResponse response;
    ClientContext context;
    context.set_wait_for_ready(true);
    if (user == "Alice") {
      context.AddMetadata("authorization", "2bd806c9");
      context.AddMetadata("membership", "premium");
    } else {
      context.AddMetadata("authorization", "81b637d8");
      context.AddMetadata("membership", "normal");
    }
    if (route != "") {
      context.AddMetadata("route", route);
    }
    std::unique_ptr<ClientReader<PriceResponse> > reader(
        stub_->WatchPrice(&context, request));
    bool first_read = true;
    while (reader->Read(&response)) {
      if (first_read) {
        auto metadata_hostname =
            context.GetServerInitialMetadata().find("hostname");
        if (metadata_hostname != context.GetServerInitialMetadata().end()) {
          std::cout << "server host: "
                    << std::string(metadata_hostname->second.data(),
                                   metadata_hostname->second.length())
                    << std::endl;
        }
        first_read = false;
      }
      std::cout << "grpc-coin price: " << response.price() << std::endl;
    }
    Status status = reader->Finish();
    if (!status.ok()) {
      std::cout << status.error_code() << ": " << status.error_message()
                << std::endl;
    }
  }

 private:
  std::unique_ptr<Stats::Stub> stub_;
};

int main(int argc, char** argv) {
  std::string command = "balance";
  std::string wallet_server = "localhost:18881";
  std::string stats_server = "localhost:18882";
  std::string user = "Alice";
  std::string route = "";
  bool affinity = false;
  bool watch = false;
  bool unary_watch = false;
  std::string gcp_client_project = "";
  std::string arg_command_balance("balance");
  std::string arg_command_price("price");
  std::string arg_str_wallet_server("--wallet_server");
  std::string arg_str_stats_server("--stats_server");
  std::string arg_str_user("--user");
  std::string arg_str_watch("--watch");
  std::string arg_str_unary_watch("--unary_watch");
  std::string arg_str_gcp_client_project("--gcp_client_project");
  std::string arg_str_route("--route");
  std::string arg_str_affinity("--affinity");
  std::string creds_type =
      traffic_director_grpc_examples::ParseCommandLineArgForCredsType(argc,
                                                                      argv);
  for (int i = 1; i < argc; ++i) {
    std::string arg_val = argv[i];
    size_t start_pos = arg_val.find(arg_command_balance);
    if (start_pos != std::string::npos) {
      command = "balance";
      continue;
    }
    start_pos = arg_val.find(arg_command_price);
    if (start_pos != std::string::npos) {
      command = "price";
      continue;
    };
    start_pos = arg_val.find(arg_str_wallet_server);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_wallet_server.size();
      if (arg_val[start_pos] == '=') {
        wallet_server = arg_val.substr(start_pos + 1);
        continue;
      } else {
        std::cout << "The only correct argument syntax is --wallet_server="
                  << std::endl;
        return 1;
      }
    }
    start_pos = arg_val.find(arg_str_stats_server);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_stats_server.size();
      if (arg_val[start_pos] == '=') {
        stats_server = arg_val.substr(start_pos + 1);
        continue;
      } else {
        std::cout << "The only correct argument syntax is --stats_server="
                  << std::endl;
        return 1;
      }
    }
    start_pos = arg_val.find(arg_str_user);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_user.size();
      if (arg_val[start_pos] == '=') {
        user = arg_val.substr(start_pos + 1);
        continue;
      } else {
        std::cout << "The only correct argument syntax is --user=" << std::endl;
        return 1;
      }
    }
    start_pos = arg_val.find(arg_str_watch);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_watch.size();
      if (arg_val[start_pos] == '=') {
        if (arg_val.substr(start_pos + 1) == "true") {
          watch = true;
          continue;
        } else if (arg_val.substr(start_pos + 1) == "false") {
          watch = false;
          continue;
        } else {
          std::cout
              << "The only correct value for argument --watch is true or false"
              << std::endl;
          return 1;
        }
      } else {
        std::cout << "The only correct argument syntax is --watch="
                  << std::endl;
        return 1;
      }
    }
    start_pos = arg_val.find(arg_str_unary_watch);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_unary_watch.size();
      if (arg_val[start_pos] == '=') {
        if (arg_val.substr(start_pos + 1) == "true") {
          if (command != "balance") {
            std::cout << "The argument --unary_watch is only applicable to "
                         "command balance"
                      << std::endl;
            return 1;
          } else if (watch == true) {
            std::cout << "The argument --unary_watch is only applicable if "
                         "--watch is set to false"
                      << std::endl;
            return 1;
          }
          unary_watch = true;
          continue;
        } else if (arg_val.substr(start_pos + 1) == "false") {
          unary_watch = false;
          continue;
        } else {
          std::cout << "The only correct value for argument --unary_watch is "
                       "true or false"
                    << std::endl;
          return 1;
        }
      } else {
        std::cout << "The only correct argument syntax is --unary_watch="
                  << std::endl;
        return 1;
      }
    }

    start_pos = arg_val.find(arg_str_gcp_client_project);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_gcp_client_project.size();
      if (arg_val[start_pos] == '=') {
        gcp_client_project = arg_val.substr(start_pos + 1);
        continue;
      } else {
        std::cout << "The only correct argument syntax is --gcp_client_project="
                  << std::endl;
        return 1;
      }
    }

    start_pos = arg_val.find(arg_str_route);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_route.size();
      if (arg_val[start_pos] == '=') {
        route = arg_val.substr(start_pos + 1);
        continue;
      } else {
        std::cout << "The only correct argument syntax is --route="
                  << std::endl;
        return 1;
      }
    }

    start_pos = arg_val.find(arg_str_affinity);
    if (start_pos != std::string::npos) {
      start_pos += arg_str_affinity.size();
      if (arg_val[start_pos] == '=') {
        if (arg_val.substr(start_pos + 1) == "true") {
          affinity = true;
          continue;
        } else if (arg_val.substr(start_pos + 1) != "false") {
          std::cout << "The only correct value for argument --affinity is "
                       "true or false"
                    << std::endl;
          return 1;
        }
      } else {
        std::cout << "The only correct argument syntax is --affinity="
                  << std::endl;
        return 1;
      }
    }
  }

  std::cout << "Client arguments: command: " << command
            << ", wallet_server: " << wallet_server
            << ", stats_server: " << stats_server << ", user: " << user
            << ", watch: " << watch << " ,unary_watch: " << unary_watch
            << ", gcp_client_project: " << gcp_client_project
            << ", route: " << route << ", affinitey: " << affinity
            << ", creds: " << creds_type << std::endl;

  if (!gcp_client_project.empty()) {
    grpc::RegisterOpenCensusPlugin();
    grpc::RegisterOpenCensusViewsForExport();
    opencensus::trace::TraceConfig::SetCurrentTraceParams(
        {128, 128, 128, 128, opencensus::trace::ProbabilitySampler(1.0)});
    opencensus::exporters::trace::StackdriverOptions trace_opts;
    trace_opts.project_id = gcp_client_project;
    opencensus::exporters::trace::StackdriverExporter::Register(
        std::move(trace_opts));
    opencensus::exporters::stats::StackdriverOptions stats_opts;
    stats_opts.project_id = gcp_client_project;
    // This must be unique among all processes exporting to Stackdriver
    stats_opts.opencensus_task = "client-" + std::to_string(getpid());
    opencensus::exporters::stats::StackdriverExporter::Register(
        std::move(stats_opts));
  }

  // Instantiate the client.  It requires a channel, out of which the actual
  // RPCs are created.  The channel models a connection to an endpoint (Stats
  // Server or Wallet Server in this case).  We indicate that the channel isn't
  // authenticated (use of InsecureChannelCredentials()).
  ChannelArguments args;
  if (command == "price") {
    StatsClient stats(grpc::CreateCustomChannel(
        stats_server,
        traffic_director_grpc_examples::GetChannelCredetials(creds_type),
        args));
    if (watch) {
      stats.WatchPrice(user, route);
    } else {
      stats.FetchPrice(user, route);
    }
  } else {
    WalletClient wallet(grpc::CreateCustomChannel(
        wallet_server,
        traffic_director_grpc_examples::GetChannelCredetials(creds_type),
        args));
    if (watch) {
      wallet.WatchBalance(user, route);
    } else {
      while (true) {
        wallet.FetchBalance(user, route, affinity);
        if (!unary_watch) break;
        std::this_thread::sleep_for(std::chrono::milliseconds(1000));
      }
    }
  }
  return 0;
}
