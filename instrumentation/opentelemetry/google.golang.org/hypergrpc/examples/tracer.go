package examples

import (
	"log"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

// InitTracer initializes the tracer and register it globally
func InitTracer(serviceName string) func() {
	// Create stdout exporter to be able to retrieve
	// the collected spans.
	stdoutExporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	// Creates a zipkin exporter that can be plugged into a hypertrace
	// ingester (see https://raw.githubusercontent.com/hypertrace/hypertrace/main/docker/docker-compose.yml)
	zipkinBatchExporter, err := zipkin.NewRawExporter("http://localhost:9411/api/v2/spans", serviceName)
	if err != nil {
		log.Fatal(err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(stdoutExporter),
		sdktrace.WithBatcher(zipkinBatchExporter, sdktrace.WithMaxExportBatchSize(1)),
		sdktrace.WithResource(resource.New(semconv.ServiceNameKey.String(serviceName))))
	if err != nil {
		log.Fatal(err)
	}

	global.SetTraceProvider(tp)
	return func() {
		<-time.After(2 * time.Second)
	}
}
