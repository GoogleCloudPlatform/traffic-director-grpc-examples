/*
 *
 * Copyright 2021 Google LLC
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

#include "utility.h"

#include <cassert>
#include <iostream>

#include "absl/memory/memory.h"
#include "grpcpp/ext/admin_services.h"
#include "grpcpp/xds_server_builder.h"

namespace traffic_director_grpc_examples {

const char* kArgCreds = "--creds";
const char* kInsecureCreds = "insecure";
const char* kXdsCreds = "xds";

const char* ParseCommandLineArgForCredsType(int argc, char** argv) {
  for (int i = 1; i < argc; ++i) {
    std::string arg = argv[i];
    size_t start_pos = arg.find(kArgCreds);
    if (start_pos != std::string::npos) {
      start_pos += strlen(kArgCreds);
      if (arg[start_pos] == '=') {
        absl::string_view arg_val = arg.substr(start_pos + 1);
        if (arg_val == kXdsCreds) {
          return kXdsCreds;
        } else if (arg_val == kInsecureCreds) {
          return kInsecureCreds;
        } else {
          std::cerr << "Allowed values for --creds are \'xds\', \'insecure\'";
          exit(1);
        }
      } else {
        std::cerr << "The only correct argument syntax is --creds=<value>";
        exit(1);
      }
    }
  }
  return kInsecureCreds;
}

std::shared_ptr<grpc::ChannelCredentials> GetChannelCredetials(
    absl::string_view creds_type) {
  if (creds_type == kInsecureCreds) {
    return grpc::InsecureChannelCredentials();
  } else if (creds_type == kXdsCreds) {
    return grpc::experimental::XdsCredentials(
        grpc::InsecureChannelCredentials());
  }
  assert(0);
}

std::unique_ptr<grpc::ServerBuilder> GetServerBuilder(
    absl::string_view creds_type) {
  if (creds_type == kInsecureCreds) {
    return absl::make_unique<grpc::ServerBuilder>();
  } else if (creds_type == kXdsCreds) {
    return absl::make_unique<grpc::experimental::XdsServerBuilder>();
  }
  assert(0);
}

std::shared_ptr<grpc::ServerCredentials> GetServerCredentials(
    absl::string_view creds_type) {
  if (creds_type == kInsecureCreds) {
    return grpc::InsecureServerCredentials();
  } else if (creds_type == kXdsCreds) {
    return grpc::experimental::XdsServerCredentials(
        grpc::InsecureServerCredentials());
  }
  assert(0);
}

std::unique_ptr<grpc::Server> StartAdminServer(const std::string& port) {
  std::string server_address = "localhost" + port;
  grpc::EnableDefaultHealthCheckService(true);
  grpc::ServerBuilder builder;
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  grpc::AddAdminServices(&builder);
  std::cout << "Admin and Health Check Server listening on " << server_address
            << std::endl;
  return builder.BuildAndStart();
}

}  // namespace traffic_director_grpc_examples
