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

import static java.util.concurrent.TimeUnit.MILLISECONDS;
import static java.util.concurrent.TimeUnit.SECONDS;

import com.google.common.util.concurrent.ListenableScheduledFuture;
import com.google.common.util.concurrent.ListeningScheduledExecutorService;
import com.google.common.util.concurrent.MoreExecutors;
import io.grpc.ChannelCredentials;
import io.grpc.examples.wallet.account.AccountGrpc;
import io.grpc.examples.wallet.account.GetUserInfoRequest;
import io.grpc.examples.wallet.account.GetUserInfoResponse;
import io.grpc.examples.wallet.account.MembershipType;
import io.grpc.examples.wallet.stats.PriceRequest;
import io.grpc.examples.wallet.stats.PriceResponse;
import io.grpc.examples.wallet.stats.StatsGrpc;
import io.grpc.Grpc;
import io.grpc.health.v1.HealthCheckResponse.ServingStatus;
import io.grpc.InsecureChannelCredentials;
import io.grpc.InsecureServerCredentials;
import io.grpc.ManagedChannel;
import io.grpc.protobuf.services.ProtoReflectionService;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.ServerCredentials;
import io.grpc.ServerInterceptors;
import io.grpc.services.AdminInterface;
import io.grpc.services.HealthStatusManager;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import io.grpc.stub.ServerCallStreamObserver;
import io.grpc.stub.StreamObserver;
import io.grpc.xds.XdsChannelCredentials;
import io.grpc.xds.XdsServerBuilder;
import io.grpc.xds.XdsServerCredentials;
import java.io.IOException;
import java.util.concurrent.Executors;
import java.util.logging.Level;
import java.util.logging.Logger;

/** Stats server for the gRPC Wallet example. */
public class StatsServer {
  private static final Logger logger = Logger.getLogger(StatsServer.class.getName());

  private enum CredentialsType {
    INSECURE,
    XDS
  }
  private Server server;
  private Server healthServer;
  private Server adminServer;

  private int port = 18882;
  private int adminPort = 28882;
  private String accountServer = "localhost:18883";
  private String hostnameSuffix = "";
  private String gcpClientProject = "";
  private boolean premiumOnly;
  private CredentialsType credentialsType = CredentialsType.INSECURE;

  private ManagedChannel accountChannel;
  private ListeningScheduledExecutorService exec;

  void parseArgs(String[] args) {
    boolean usage = false;
    for (String arg : args) {
      if (!arg.startsWith("--")) {
        System.err.println("All arguments must start with '--': " + arg);
        usage = true;
        break;
      }
      String[] parts = arg.substring(2).split("=", 2);
      String key = parts[0];
      if ("help".equals(key)) {
        usage = true;
        break;
      }
      if (parts.length != 2) {
        System.err.println("All flags must be of the form --arg=value");
        usage = true;
        break;
      }
      String value = parts[1];
      if ("port".equals(key)) {
        port = Integer.parseInt(value);
      } else if ("admin_port".equals(key)) {
        adminPort = Integer.parseInt(value);
      } else if ("account_server".equals(key)) {
        accountServer = value;
      } else if ("hostname_suffix".equals(key)) {
        hostnameSuffix = value;
      } else if ("gcp_client_project".equals(key)) {
        gcpClientProject = value;
      } else if ("premium_only".equals(key)) {
        premiumOnly = Boolean.parseBoolean(value);
      } else if ("creds".equals(key)) {
        credentialsType = CredentialsType.valueOf(value.toUpperCase());
      } else {
        System.err.println("Unknown argument: " + key);
        usage = true;
        break;
      }
    }
    if (usage) {
      StatsServer s = new StatsServer();
      System.out.println(
          "Usage: [ARGS...]"
              + "\n"
              + "\n  --port=PORT                The port to listen on. Default "
              + s.port
              + "\n  --admin_port=PORT          The admin port to listen on. Default "
              + s.adminPort
              + "\n  --account_server=HOST      Address of the account server. Default "
              + s.accountServer
              + "\n  --hostname_suffix=STR      Suffix to append to hostname in response header. "
              + "Default \""
              + s.hostnameSuffix
              + "\""
              + "\n  --gcp_client_project=STR GCP project. If set, metrics and traces will be "
              + "sent to Stackdriver. Default \"" + s.gcpClientProject + "\""
              + "\n  --premium_only=true|false  If true, all non-premium RPCs are rejected. "
              + "Default "
              + s.premiumOnly
              + "\n  --creds=insecure|xds  . Type of credentials to use on the client & server. "
              + "Default "
              + s.credentialsType.toString().toLowerCase());
      System.exit(1);
    }
  }

  private void start() throws IOException {
    if (!gcpClientProject.isEmpty()) {
      Observability.registerExporters(gcpClientProject);
    }
    adminServer = ServerBuilder.forPort(adminPort)
        .addServices(AdminInterface.getStandardServices())
        .build()
        .start();
    logger.info("Admin server started, listening on " + adminPort);
    ChannelCredentials channelCredentials =
        credentialsType == CredentialsType.XDS
            ? XdsChannelCredentials.create(InsecureChannelCredentials.create())
            : InsecureChannelCredentials.create();
    accountChannel = Grpc.newChannelBuilder(accountServer, channelCredentials).build();
    exec = MoreExecutors.listeningDecorator(Executors.newSingleThreadScheduledExecutor());
    HealthStatusManager health = new HealthStatusManager();
    ServerCredentials serverCredentials =
        credentialsType == CredentialsType.XDS
            ? XdsServerCredentials.create(InsecureServerCredentials.create())
            : InsecureServerCredentials.create();
    // Since the main server may be using TLS, we start a second server just for plaintext health
    // checks
    int healthPort = port + 1;
    if (credentialsType == CredentialsType.XDS) {
      server =
          XdsServerBuilder.forPort(port, serverCredentials)
              .addService(
                  ServerInterceptors.intercept(
                      new StatsImpl(accountChannel, exec, premiumOnly),
                      new WalletInterceptors.HostnameInterceptor(),
                      new WalletInterceptors.AuthInterceptor()))
              .addService(ProtoReflectionService.newInstance())
              .addService(health.getHealthService())
              .build()
              .start();
      healthServer =
          XdsServerBuilder.forPort(healthPort, InsecureServerCredentials.create())
              .addService(health.getHealthService()) // allow management servers to monitor health
              .build()
              .start();
    } else {
      server =
          ServerBuilder.forPort(port)
              .addService(
                  ServerInterceptors.intercept(
                      new StatsImpl(accountChannel, exec, premiumOnly),
                      new WalletInterceptors.HostnameInterceptor(),
                      new WalletInterceptors.AuthInterceptor()))
              .addService(ProtoReflectionService.newInstance())
              .addService(health.getHealthService())
              .build()
              .start();
      healthServer =
          ServerBuilder.forPort(healthPort)
              .addService(health.getHealthService()) // allow management servers to monitor health
              .build()
              .start();
    }
    health.setStatus("", ServingStatus.SERVING);
    logger.info("Server started, listening on " + port);
    logger.info("Plaintext health server started, listening on " + healthPort);
    Runtime.getRuntime()
        .addShutdownHook(
            new Thread() {
              @Override
              public void run() {
                System.err.println("*** shutting down gRPC server since JVM is shutting down");
                try {
                  StatsServer.this.stop();
                } catch (InterruptedException e) {
                  e.printStackTrace(System.err);
                }
                System.err.println("*** server shut down");
              }
            });
  }

  private void stop() throws InterruptedException {
    if (server != null) {
      server.shutdown().awaitTermination(30, SECONDS);
    }
    if (healthServer != null) {
      healthServer.shutdown().awaitTermination(30, SECONDS);
    }
    if (adminServer != null) {
      adminServer.shutdown().awaitTermination(30, SECONDS);
    }
    if (accountChannel != null) {
      accountChannel.shutdownNow().awaitTermination(5, SECONDS);
    }
    if (exec != null) {
      exec.shutdownNow();
    }
  }

  private void blockUntilShutdown() throws InterruptedException {
    if (server != null) {
      server.awaitTermination();
    }
    if (healthServer != null) {
      healthServer.awaitTermination();
    }
  }

  public static void main(String[] args) throws IOException, InterruptedException {
    final StatsServer server = new StatsServer();
    server.parseArgs(args);
    server.start();
    server.blockUntilShutdown();
  }

  private static class StatsImpl extends StatsGrpc.StatsImplBase {
    private final AccountGrpc.AccountBlockingStub blockingStub;
    private final ListeningScheduledExecutorService exec;
    private final boolean premiumOnly;

    private StatsImpl(
        ManagedChannel accountChannel,
        ListeningScheduledExecutorService exec,
        boolean premiumOnly) {
      this.blockingStub = AccountGrpc.newBlockingStub(accountChannel);
      this.exec = exec;
      this.premiumOnly = premiumOnly;
    }

    private boolean validateMembership(
        String token, String membership, StreamObserver<PriceResponse> responseObserver) {
      try {
        GetUserInfoResponse response =
            blockingStub.getUserInfo(GetUserInfoRequest.newBuilder().setToken(token).build());
        MembershipType type = response.getMembership();
        if ("premium".equals(membership) && type != MembershipType.PREMIUM) {
          responseObserver.onError(
              Status.UNAUTHENTICATED
                  .withDescription("token does not belong to a premium member")
                  .asRuntimeException());
          return false;
        }
        if (premiumOnly && !"premium".equals(membership)) {
          responseObserver.onError(
              Status.PERMISSION_DENIED
                  .withDescription("only premium RPCs are allowed by this service")
                  .asRuntimeException());
          return false;
        }
      } catch (StatusRuntimeException e) {
        logger.log(Level.WARNING, "RPC failed: {0}", e.getStatus());
        responseObserver.onError(
            Status.INTERNAL
                .withDescription("Failed to connect to account server")
                .asRuntimeException());
        return false;
      }
      return true;
    }

    private long getPrice() {
      return Double.valueOf(Math.sin(System.currentTimeMillis() / 173) * 1000 + 10000).longValue();
    }

    @Override
    public void watchPrice(PriceRequest req, final StreamObserver<PriceResponse> responseObserver) {
      String token = WalletInterceptors.TOKEN_KEY.get();
      String membership = WalletInterceptors.MEMBERSHIP_KEY.get();
      if (!validateMembership(token, membership, responseObserver)) {
        return;
      }
      int millisecondsBetweenUpdates;
      if ("premium".equals(membership)) {
        millisecondsBetweenUpdates = 100;
      } else {
        millisecondsBetweenUpdates = 1000;
      }
      final ListenableScheduledFuture<?> future =
          exec.scheduleAtFixedRate(
              new Runnable() {
                @Override
                public void run() {
                  PriceResponse response = PriceResponse.newBuilder().setPrice(getPrice()).build();
                  responseObserver.onNext(response);
                }
              },
              0,
              millisecondsBetweenUpdates,
              MILLISECONDS);

      ((ServerCallStreamObserver) responseObserver)
          .setOnCancelHandler(
              new Runnable() {
                @Override
                public void run() {
                  future.cancel(true);
                }
              });
    }

    @Override
    public void fetchPrice(PriceRequest req, StreamObserver<PriceResponse> responseObserver) {
      String token = WalletInterceptors.TOKEN_KEY.get();
      String membership = WalletInterceptors.MEMBERSHIP_KEY.get();

      if (!validateMembership(token, membership, responseObserver)) {
        return;
      }

      PriceResponse response = PriceResponse.newBuilder().setPrice(getPrice()).build();
      responseObserver.onNext(response);
      responseObserver.onCompleted();
    }
  }
}
