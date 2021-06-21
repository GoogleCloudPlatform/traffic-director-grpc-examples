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

#include <memory>

#include "absl/strings/string_view.h"
#include "grpcpp/grpcpp.h"
#include "grpcpp/security/credentials.h"
#include "grpcpp/security/server_credentials.h"

namespace traffic_director_grpc_examples {

extern const char* kArgCreds;
extern const char* kInsecureCreds;
extern const char* kXdsCreds;

// Parses command-line args for --creds=<value>
// Allowed values are 'xds' and 'insecure'.
// Returns 'insecure' if no matching arg is found.
const char* ParseCommandLineArgForCredsType(int argc, char** argv);

std::shared_ptr<grpc::ChannelCredentials> GetChannelCredetials(
    absl::string_view creds_type);

std::unique_ptr<grpc::ServerBuilder> GetServerBuilder(
    absl::string_view creds_type);

std::shared_ptr<grpc::ServerCredentials> GetServerCredentials(
    absl::string_view creds_type);

// Starts admin server on port \a port using insecure credentials.
// This server also exposes a health-check service on the same port.
std::unique_ptr<grpc::Server> StartAdminServer(const std::string& port);

}  // namespace traffic_director_grpc_examples
