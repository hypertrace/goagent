package otel

import (
	"log"

	"github.com/traceableai/goagent"
	grpcserver "github.com/traceableai/goagent/otel/grpc/server"
	httpserver "github.com/traceableai/goagent/otel/http/server"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

// init loads implicitly the instrumentation elements and initializes the tracer
// TODO: Define settings for the tracer
func init() {
	initTracer()
	goagent.Instrumentation.HTTPHandler = httpserver.NewHandler
	goagent.Instrumentation.GRPCInterceptor.UnaryServer = grpcserver.NewUnaryServerInterceptor
}

func initTracer() {
	// Create stdout exporter to be able to retrieve
	// the collected spans.
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.New(semconv.ServiceNameKey.String("ExampleService"))))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
}
