package observability

import (
	"log"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func ConfigureStackdriver(gcpProject string) {
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatalf("Failed to register ocgrpc client views: %v", err)
	}

	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: gcpProject,
	})
	if err != nil {
		log.Fatalf("Failed to create Stackdriver exporter: %v", err)
	}
	defer sd.Flush()
	trace.RegisterExporter(sd)
	sd.StartMetricsExporter()
	defer sd.StopMetricsExporter()
	// For demo purposes, always sample
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}
