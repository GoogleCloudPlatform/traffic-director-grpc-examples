package io.grpc.examples.wallet;

import static com.google.common.truth.Truth.assertThat;

import com.google.common.collect.ImmutableSet;
import io.grpc.ManagedChannel;
import io.grpc.Server;
import io.grpc.netty.shaded.io.grpc.netty.NettyChannelBuilder;
import io.grpc.netty.shaded.io.grpc.netty.NettyServerBuilder;
import io.grpc.stub.StreamObserver;
import io.opencensus.common.Duration;
import io.opencensus.contrib.grpc.metrics.RpcViews;
import io.opencensus.exporter.metrics.util.IntervalMetricReader;
import io.opencensus.exporter.metrics.util.MetricExporter;
import io.opencensus.exporter.metrics.util.MetricReader;
import io.opencensus.metrics.Metrics;
import io.opencensus.metrics.export.Metric;
import io.opencensus.metrics.export.MetricDescriptor;
import io.opencensus.trace.Tracing;
import io.opencensus.trace.config.TraceConfig;
import io.opencensus.trace.export.SpanData;
import io.opencensus.trace.export.SpanExporter;
import io.opencensus.trace.samplers.Samplers;
import java.util.Collection;
import java.util.HashSet;
import java.util.Set;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;

@RunWith(JUnit4.class)
public class ObservabilityTest {

  private static final int EXPORT_WAIT_SEC = 15;

  private static final Set<String> METRICS =
      ImmutableSet.of(
          "grpc.io/client/roundtrip_latency",
          "grpc.io/client/completed_rpcs",
          "grpc.io/client/sent_bytes_per_rpc",
          "grpc.io/client/received_bytes_per_rpc",
          "grpc.io/client/sent_messages_per_rpc",
          "grpc.io/client/received_messages_per_rpc",
          "grpc.io/server/completed_rpcs",
          "grpc.io/server/sent_bytes_per_rpc",
          "grpc.io/server/received_bytes_per_rpc",
          "grpc.io/server/sent_messages_per_rpc",
          "grpc.io/server/received_messages_per_rpc",
          "grpc.io/server/server_latency");

  private Server server;
  private ManagedChannel channel;
  private RecordingMetricsExporter recordingMetricsExporter;
  private RecordingTraceExporter recordingTraceExporter;

  @Before
  public void setUp() throws Exception {
    RpcViews.registerAllGrpcViews();
    recordingMetricsExporter = new RecordingMetricsExporter();

    TraceConfig traceConfig = Tracing.getTraceConfig();
    traceConfig.updateActiveTraceParams(
        traceConfig.getActiveTraceParams().toBuilder().setSampler(Samplers.alwaysSample()).build());
    recordingTraceExporter = new RecordingTraceExporter();
    Tracing.getExportComponent()
        .getSpanExporter()
        .registerHandler(RecordingTraceExporter.class.getName(), recordingTraceExporter);

    server = NettyServerBuilder.forPort(0).addService(new WalletImpl()).build();
    server.start();
    channel = NettyChannelBuilder.forAddress("localhost", server.getPort()).usePlaintext().build();
  }

  @After
  public void tearDown() {
    channel.shutdownNow();
    server.shutdownNow();
  }

  @Test
  public void testMetrics() throws Exception {
    WalletGrpc.newBlockingStub(channel).fetchBalance(BalanceRequest.getDefaultInstance());
    recordingMetricsExporter.latch.await(EXPORT_WAIT_SEC, TimeUnit.SECONDS);
    assertThat(recordingMetricsExporter.metrics).containsAtLeastElementsIn(METRICS);
  }

  @Test
  public void testTrace() throws Exception {
    WalletGrpc.newBlockingStub(channel).fetchBalance(BalanceRequest.getDefaultInstance());
    recordingTraceExporter.latch.await(EXPORT_WAIT_SEC, TimeUnit.SECONDS);
    assertThat(recordingTraceExporter.clientSpanData).isNotNull();
    assertThat(recordingTraceExporter.serverSpanData).isNotNull();
    assertThat(recordingTraceExporter.serverSpanData.getContext().getTraceId())
        .isEqualTo(recordingTraceExporter.clientSpanData.getContext().getTraceId());
    assertThat(recordingTraceExporter.serverSpanData.getParentSpanId())
        .isEqualTo(recordingTraceExporter.clientSpanData.getContext().getSpanId());
  }

  private static class WalletImpl extends WalletGrpc.WalletImplBase {
    @Override
    public void fetchBalance(
        BalanceRequest request, StreamObserver<BalanceResponse> responseObserver) {
      responseObserver.onNext(BalanceResponse.getDefaultInstance());
      responseObserver.onCompleted();
    }
  }

  private class RecordingTraceExporter extends SpanExporter.Handler {
    // The trace span exported by the client (identified by a null parent span)
    SpanData clientSpanData = null;
    // The trace span exported by the server (identified by a non-null parent span)
    SpanData serverSpanData = null;
    CountDownLatch latch = new CountDownLatch(2);

    @Override
    public void export(Collection<SpanData> spanDataList) {
      for (SpanData sd : spanDataList) {
        if (sd.getParentSpanId() != null) {
          assertThat(serverSpanData).isNull();
          serverSpanData = sd;
        } else {
          assertThat(clientSpanData).isNull();
          clientSpanData = sd;
        }
        latch.countDown();
      }
    }
  }

  private class RecordingMetricsExporter extends MetricExporter {
    private final Set<String> metrics = new HashSet<>();
    private final CountDownLatch latch = new CountDownLatch(1);

    private RecordingMetricsExporter() {
      IntervalMetricReader.Options.Builder options = IntervalMetricReader.Options.builder();
      MetricReader reader =
          MetricReader.create(
              MetricReader.Options.builder()
                  .setMetricProducerManager(Metrics.getExportComponent().getMetricProducerManager())
                  .build());
      IntervalMetricReader.create(
          this, reader, options.setExportInterval(Duration.create(1, 0)).build());
    }

    @Override
    public void export(Collection<Metric> metricsToExport) {
      for (Metric metric : metricsToExport) {
        if (!metric.getTimeSeriesList().isEmpty()) {
          MetricDescriptor md = metric.getMetricDescriptor();
          metrics.add(md.getName());
        }
      }
      if (metrics.containsAll(METRICS)) {
        latch.countDown();
      }
    }
  }
}
