package examples

import (
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/trace"
)

// InitTracer initializes the tracer and register it globally
func InitTracer(serviceName string) {
	// Register stats and trace exporters to export the collected data.
	exporter := &exporter.PrintExporter{}
	trace.RegisterExporter(exporter)
	// Always trace for this demo. In a production application, you should
	// configure this to a trace.ProbabilitySampler set at the desired
	// probability.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}
