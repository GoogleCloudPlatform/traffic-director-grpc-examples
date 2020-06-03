/*
 * Copyright 2020 The gRPC Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package io.grpc.examples.wallet;

import static io.grpc.Metadata.ASCII_STRING_MARSHALLER;

import io.grpc.Context;
import io.grpc.Contexts;
import io.grpc.ForwardingServerCall.SimpleForwardingServerCall;
import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Random;
import java.util.logging.Level;
import java.util.logging.Logger;

/**
 * Interceptor to add hostname to outbound metadata and read token and membership type into context.
 */
final class WalletServerInterceptor implements ServerInterceptor {
  private static final Logger logger = Logger.getLogger(WalletServerInterceptor.class.getName());

  public static final Context.Key<String> TOKEN_KEY = Context.key("token");
  public static final Context.Key<String> PREMIUM_KEY = Context.key("premium");
  public static final Metadata.Key<String> TOKEN_MD_KEY =
      Metadata.Key.of("Authorization", ASCII_STRING_MARSHALLER);
  public static final Metadata.Key<String> PREMIUM_MD_KEY =
      Metadata.Key.of("premium", ASCII_STRING_MARSHALLER);
  public static final Metadata.Key<String> HOSTNAME_MD_KEY =
      Metadata.Key.of("hostname", ASCII_STRING_MARSHALLER);

  private final String hostname;

  WalletServerInterceptor() {
    String hostname;
    try {
      hostname = InetAddress.getLocalHost().getHostName();
    } catch (UnknownHostException e) {
      logger.log(Level.WARNING, "Failed to get host", e);
      Random random = new Random();
      hostname = String.format("generated-%3d", random.nextInt(1000));
      logger.log(Level.WARNING, "Using " + hostname + " as hostname");
    }
    this.hostname = hostname;
  }

  @Override
  public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
      ServerCall<ReqT, RespT> call, Metadata requestHeaders, ServerCallHandler<ReqT, RespT> next) {
    Context newContext =
        Context.current()
            .withValue(TOKEN_KEY, requestHeaders.get(TOKEN_MD_KEY))
            .withValue(PREMIUM_KEY, requestHeaders.get(PREMIUM_MD_KEY));
    return Contexts.interceptCall(
        newContext,
        new SimpleForwardingServerCall<ReqT, RespT>(call) {
          @Override
          public void sendHeaders(Metadata responseHeaders) {
            responseHeaders.put(HOSTNAME_MD_KEY, hostname);
            super.sendHeaders(responseHeaders);
          }
        },
        requestHeaders,
        next);
  }
}
