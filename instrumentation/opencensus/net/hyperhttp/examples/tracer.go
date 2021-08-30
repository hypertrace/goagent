package examples // import "github.com/hypertrace/goagent/instrumentation/opencensus/net/hyperhttp/examples"

import (
	"contrib.go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/trace"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

// InitTracer initializes the tracer and register it globally
func InitTracer(serviceName string) func() error {
	// Register stats and trace exporters to export the collected data.
	stdoutExporter := &exporter.PrintExporter{}
	trace.RegisterExporter(stdoutExporter)

	// Creates a zipkin exporter that can be plugged into a hypertrace
	// ingester (see https://raw.githubusercontent.com/hypertrace/hypertrace/main/docker/docker-compose.yml)
	reporterURI := "http://localhost:9411/api/v2/spans"
	localEndpoint, _ := openzipkin.NewEndpoint(serviceName, "localhost")
	reporter := zipkinHTTP.NewReporter(reporterURI)
	zipkinExporter := zipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(zipkinExporter)

	// Always trace for this demo. In a production application, you should
	// configure this to a trace.ProbabilitySampler set at the desired
	// probability.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return reporter.Close
}
