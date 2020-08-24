package internal

import (
	"log"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/standard"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitTracer initializes the tracer and returns a flusher of the reported
// span for further inspection
func InitTracer() func() []*trace.SpanData {
	exporter := &Recorder{}

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.New(standard.ServiceNameKey.String("ExampleService"))))
	if err != nil {
		log.Fatal(err)
	}

	global.SetTraceProvider(tp)
	global.SetPropagators(propagation.New(
		propagation.WithExtractors(apitrace.B3{}),
	))

	return func() []*trace.SpanData {
		return exporter.Flush()
	}
}
