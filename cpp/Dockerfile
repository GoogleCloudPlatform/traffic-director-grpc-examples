# Copyright 2021 The gRPC Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Dockerfile for building an image with all binaries required for the C++ wallet example.
# To build the image, run the following command (after replacing <TAG>) from the traffic-director-grpc-examples directory:
# docker build -t <TAG> -f cpp/Dockerfile .

FROM phusion/baseimage:master@sha256:74f8b98541d539563be2a21eefbe4b641ad43b779880b76fc02ea87b7b2ce489

RUN apt-get update -y && \
        apt-get install -y \
            build-essential \
            clang \
            python3 \
            python3-dev \
            apt-transport-https \
            curl \
            gnupg
RUN curl -fsSL https://bazel.build/bazel-release.pub.gpg | gpg --dearmor > bazel.gpg
RUN mv bazel.gpg /etc/apt/trusted.gpg.d/
RUN echo "deb [arch=amd64] https://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list
RUN apt-get update -y && apt-get install -y bazel

WORKDIR /workdir

RUN ln -s /usr/bin/python3 /usr/bin/python
RUN mkdir /artifacts

COPY . .
RUN bazel build //cpp:all
RUN cp /workdir/bazel-bin/cpp/client /artifacts/
RUN cp /workdir/bazel-bin/cpp/wallet-server /artifacts/
RUN cp /workdir/bazel-bin/cpp/account-server /artifacts/
RUN cp /workdir/bazel-bin/cpp/stats-server /artifacts/

FROM phusion/baseimage:master@sha256:74f8b98541d539563be2a21eefbe4b641ad43b779880b76fc02ea87b7b2ce489
COPY --from=0 /artifacts ./

