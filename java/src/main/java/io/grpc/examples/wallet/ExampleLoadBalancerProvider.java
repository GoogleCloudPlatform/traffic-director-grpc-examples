/*
 * Copyright 2022 Google LLC
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

import com.google.common.collect.ImmutableMap;
import io.grpc.LoadBalancer;
import io.grpc.LoadBalancer.Helper;
import io.grpc.LoadBalancerProvider;
import io.grpc.LoadBalancerRegistry;
import io.grpc.NameResolver;
import io.grpc.NameResolver.ConfigOrError;
import io.grpc.Status;
import io.grpc.util.ForwardingLoadBalancer;
import java.util.Map;


/**
 * This {@link LoadBalancerProvider} provides an example {@link LoadBalancer} implementation that
 * delegates to round_robin and simply prints to System.out when it needs to handle resolved
 * addresses. The purpose of this implementation is to demonstrate how to configure a custom {@link
 * LoadBalancer} in Traffic Director, as explained in:
 * https://cloud.google.com/traffic-director/docs/proxyless-configure-advanced-traffic-management
 */
public class ExampleLoadBalancerProvider extends LoadBalancerProvider {

  @Override
  public NameResolver.ConfigOrError parseLoadBalancingPolicyConfig(
      Map<String, ?> rawLoadBalancingPolicyConfig) {
    ConfigOrError response = null;
    try {
      LoadBalancerProvider roundRobinProvider = LoadBalancerRegistry.getDefaultRegistry()
          .getProvider("round_robin");
      ConfigOrError roundRobinConfig = roundRobinProvider.parseLoadBalancingPolicyConfig(
          ImmutableMap.of());

      // The configuration map should have a "message" entry with a message that is printed out
      // every time this load balancer handles resolved addresses.
      String message = (String) rawLoadBalancingPolicyConfig.get("message");
      if (message == null) {
        response = NameResolver.ConfigOrError.fromError(
            Status.UNAVAILABLE.withDescription("no 'message' defined"));
      }
      response = NameResolver.ConfigOrError.fromConfig(
          new ExampleLoadBalancerConfig(message, roundRobinProvider, roundRobinConfig));
    } catch (RuntimeException e) {
      response = NameResolver.ConfigOrError.fromError(
          Status.UNAVAILABLE.withDescription("Failed to parse example LB service config")
              .withCause(e));
    }
    return response;
  }

  /**
   * This is the name the load balancer is registered in the gRPC load balancer registry in the
   * {@code resources/services/io.grpc.LoadBalancerProvider} file.
   */
  @Override
  public String getPolicyName() {
    return "example.ExampleLoadBalancer";
  }

  @Override
  public boolean isAvailable() {
    return true;
  }

  @Override
  public int getPriority() {
    return 5;
  }

  @Override
  public LoadBalancer newLoadBalancer(Helper helper) {
    return new ExampleLoadBalancer(helper);
  }

  /**
   * Parsed configuration for {@link ExampleLoadBalancer}.
   */
  static class ExampleLoadBalancerConfig {

    final String message;
    final LoadBalancerProvider roundRobinProvider;
    final ConfigOrError roundRobinConfig;

    ExampleLoadBalancerConfig(String message, LoadBalancerProvider roundRobinProvider,
        ConfigOrError roundRobinConfig) {
      this.message = message;
      this.roundRobinProvider = roundRobinProvider;
      this.roundRobinConfig = roundRobinConfig;
    }
  }

  /**
   * This example {@code LoadBalancer} simply forwards to another one and prints a custom message
   * each time newly resolved addresses are handled.
   */
  static class ExampleLoadBalancer extends ForwardingLoadBalancer {

    private final Helper helper;
    private LoadBalancer delegateLb;

    ExampleLoadBalancer(Helper helper) {
      this.helper = helper;
    }

    @Override
    protected LoadBalancer delegate() {
      return delegateLb;
    }

    @Override
    public void handleResolvedAddresses(ResolvedAddresses resolvedAddresses) {
      ExampleLoadBalancerConfig config
          = (ExampleLoadBalancerConfig) resolvedAddresses.getLoadBalancingPolicyConfig();
      this.delegateLb = config.roundRobinProvider.newLoadBalancer(helper);
      System.out.println(
          "ExampleLoadBalancer handling resolved addresses [message: '" + config.message + "']");
      delegateLb.handleResolvedAddresses(resolvedAddresses.toBuilder()
          .setLoadBalancingPolicyConfig(config.roundRobinConfig).build());
    }
  }
}
