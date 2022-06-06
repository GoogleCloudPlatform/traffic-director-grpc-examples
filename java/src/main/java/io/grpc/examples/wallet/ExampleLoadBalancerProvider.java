package io.grpc.examples.wallet;

import io.grpc.LoadBalancer;
import io.grpc.LoadBalancer.Helper;
import io.grpc.LoadBalancerProvider;
import io.grpc.LoadBalancerRegistry;
import io.grpc.NameResolver;
import io.grpc.NameResolver.ConfigOrError;
import io.grpc.Status;
import io.grpc.util.ForwardingLoadBalancer;
import java.util.Map;
import java.util.logging.Logger;


/**
 * This {@link LoadBalancerProvider} provides an example {@link LoadBalancer} implementation that
 * delegates to round_robin and simply prints to System.out when it needs to handle resolved
 * addresses. The purpose of this implementation is to demonstrate how to configure a custom
 * {@link LoadBalancer} in Traffic Director, as explained in:
 * https://cloud.google.com/traffic-director/docs/proxyless-configure-advanced-traffic-management
 */
public class ExampleLoadBalancerProvider extends LoadBalancerProvider {

  @Override
  public NameResolver.ConfigOrError parseLoadBalancingPolicyConfig(
      Map<String, ?> rawLoadBalancingPolicyConfig) {
    // The configuration map should have a "message" entry with a message that is printed out every
    // time this load balancer handles resolved addresses.
    String message = (String) rawLoadBalancingPolicyConfig.get("message");
    if (message == null) {
      return NameResolver.ConfigOrError.fromError(
          Status.INVALID_ARGUMENT.withDescription("no 'message' defined"));
    }

    return NameResolver.ConfigOrError.fromConfig(new ExampleLoadBalancerConfig(message));
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
    return new ExampleLoadBalancer(helper,
        LoadBalancerRegistry.getDefaultRegistry().getProvider("round_robin")
            .newLoadBalancer(helper));
  }

  /**
   * Parsed configuration for {@link ExampleLoadBalancer}.
   */
  static class ExampleLoadBalancerConfig {

    final String message;

    ExampleLoadBalancerConfig(String message) {
      this.message = message;
    }
  }

  /**
   * This example {@code LoadBalancer} simply forwards to another one and prints a custom message
   * each time newly resolved addresses are handled.
   */
  static class ExampleLoadBalancer extends ForwardingLoadBalancer {

    private final LoadBalancer delegateLb;

    ExampleLoadBalancer(Helper helper, LoadBalancer delegateLb) {
      this.delegateLb = delegateLb;
    }

    @Override
    protected LoadBalancer delegate() {
      return delegateLb;
    }

    @Override
    public void handleResolvedAddresses(ResolvedAddresses resolvedAddresses) {
      ExampleLoadBalancerConfig config
          = (ExampleLoadBalancerConfig) resolvedAddresses.getLoadBalancingPolicyConfig();
      System.out.println(
          "ExampleLoadBalancer handling resolved addresses [message: '" + config.message + "']");
      delegateLb.handleResolvedAddresses(resolvedAddresses.toBuilder()
          .setLoadBalancingPolicyConfig(ConfigOrError.fromConfig("no config")).build());
    }
  }
}
