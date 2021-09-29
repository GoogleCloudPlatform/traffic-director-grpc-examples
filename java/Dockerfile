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
# docker build -t <TAG> -f java/Dockerfile .


FROM openjdk:11

WORKDIR /workdir

COPY . .

WORKDIR /workdir/java

RUN ./gradlew installDist

RUN cp -r /workdir/java/build/install/wallet /artifacts

FROM openjdk:11

RUN mkdir -p /build/install/wallet

COPY --from=0 /artifacts /build/install/wallet

RUN ln -s /build/install/wallet/bin/account-server /account-server

RUN ln -s /build/install/wallet/bin/stats-server /stats-server

RUN ln -s /build/install/wallet/bin/wallet-server /wallet-server

RUN ln -s /build/install/wallet/bin/client /wallet-client

CMD ["/bin/sleep","inf"]

EXPOSE 8000
