package internal

import (
	"context"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	apitrace "go.opentelemetry.io/otel/trace"
)

// InitTracer initializes the tracer and returns a flusher of the reported
// spans for further inspection. Its main purpose is to declare a tracer
// for TESTING.
func InitTracer() (apitrace.Tracer, func() []*trace.SpanData) {
	exporter := &Recorder{}

	resources, _ := resource.New(context.Background(), resource.WithAttributes(semconv.ServiceNameKey.String("TestService")))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resources),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.B3{})

	return tp.Tracer(opentelemetry.TracerDomain), func() []*trace.SpanData {
		return exporter.Flush()
	}
}
