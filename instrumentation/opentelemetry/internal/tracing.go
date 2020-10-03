package internal

import (
	"log"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

// InitTracer initializes the tracer and returns a flusher of the reported
// spans for further inspection. Its main purpose is to declare a tracer
// for TESTING.
func InitTracer() (apitrace.Tracer, func() []*trace.SpanData) {
	exporter := &Recorder{}

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.New(semconv.ServiceNameKey.String("TestService"))))
	if err != nil {
		log.Fatal(err)
	}

	global.SetTraceProvider(tp)
	global.SetPropagators(propagation.New(
		propagation.WithExtractors(apitrace.B3{}),
		propagation.WithInjectors(apitrace.B3{}),
	))

	return tp.Tracer("ai.traceable"), func() []*trace.SpanData {
		return exporter.Flush()
	}
}
