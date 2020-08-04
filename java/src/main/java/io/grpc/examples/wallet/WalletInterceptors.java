/*
 * Copyright 2020 Google LLC
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
import io.grpc.Status;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Random;
import java.util.logging.Level;
import java.util.logging.Logger;

/** Implementations of server interceptors for the gRPC Wallet example. */
final class WalletInterceptors {
  private static final Logger logger = Logger.getLogger(WalletInterceptors.class.getName());

  public static final Context.Key<String> TOKEN_KEY = Context.key("token");
  public static final Context.Key<String> MEMBERSHIP_KEY = Context.key("membership");
  public static final Metadata.Key<String> TOKEN_MD_KEY =
      Metadata.Key.of("Authorization", ASCII_STRING_MARSHALLER);
  public static final Metadata.Key<String> MEMBERSHIP_MD_KEY =
      Metadata.Key.of("membership", ASCII_STRING_MARSHALLER);
  public static final Metadata.Key<String> HOSTNAME_MD_KEY =
      Metadata.Key.of("hostname", ASCII_STRING_MARSHALLER);

  /** Adds the server hostname to response metadata. */
  static class HostnameInterceptor implements ServerInterceptor {
    private final String hostname;

    HostnameInterceptor() {
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
        ServerCall<ReqT, RespT> call,
        Metadata requestHeaders,
        ServerCallHandler<ReqT, RespT> next) {
      return next.startCall(
          new SimpleForwardingServerCall<ReqT, RespT>(call) {
            @Override
            public void sendHeaders(Metadata responseHeaders) {
              responseHeaders.put(HOSTNAME_MD_KEY, hostname);
              super.sendHeaders(responseHeaders);
            }
          },
          requestHeaders);
    }
  }

  /** Extracts token and membership type from metadata and inserts into context. */
  static class AuthInterceptor implements ServerInterceptor {
    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
        ServerCall<ReqT, RespT> call,
        Metadata requestHeaders,
        ServerCallHandler<ReqT, RespT> next) {
      String token = requestHeaders.get(TOKEN_MD_KEY);
      if (token == null) {
        call.close(Status.UNAUTHENTICATED.withDescription("missing token"), new Metadata());
        return new ServerCall.Listener() {};
      }
      String membership = requestHeaders.get(MEMBERSHIP_MD_KEY);
      if (membership == null) {
        call.close(Status.UNAUTHENTICATED.withDescription("missing membership"), new Metadata());
        return new ServerCall.Listener() {};
      }
      if (!"premium".equals(membership) && !"normal".equals(membership)) {
        call.close(
            Status.UNAUTHENTICATED.withDescription("membership must be premium or normal"),
            new Metadata());
        return new ServerCall.Listener() {};
      }
      Context newContext =
          Context.current().withValue(TOKEN_KEY, token).withValue(MEMBERSHIP_KEY, membership);
      return Contexts.interceptCall(newContext, call, requestHeaders, next);
    }
  }
}
