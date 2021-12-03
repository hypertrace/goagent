package examples

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/config"
	hyperotel "github.com/hypertrace/goagent/instrumentation/opentelemetry"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/net/hyperhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var otherSpanExporter trace.SpanExporter = nil

func ExampleInitAdditional() {
	hyperSpanProcessor, shutdown := hyperotel.InitAsAdditional(config.Load())
	defer shutdown()

	ctx := context.Background()
	resources, _ := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("my-server"),
		),
	)

	otherSpanProcessor := sdktrace.NewBatchSpanProcessor(
		hyperotel.RemoveGoAgentAttrs(otherSpanExporter),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(hyperSpanProcessor),
		sdktrace.WithSpanProcessor(otherSpanProcessor),
		sdktrace.WithResource(resources),
	)

	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	r := mux.NewRouter()
	r.Handle("/foo", otelhttp.NewHandler(
		hyperhttp.WrapHandler(
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}),
			nil,
		),
		"/foo",
	))

	log.Fatal(http.ListenAndServe(":8081", r))
}
