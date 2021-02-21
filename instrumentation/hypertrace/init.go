package hypertrace

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/hypertrace/goagent/config"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	"github.com/hypertrace/goagent/version"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

const batchTimeoutInSecs = 200.0

func makePropagator(formats []config.PropagationFormat) propagation.TextMapPropagator {
	var propagators []propagation.TextMapPropagator
	for _, format := range formats {
		switch format {
		case config.PropagationFormat_B3:
			propagators = append(propagators, b3.B3{})
		case config.PropagationFormat_TRACECONTEXT:
			propagators = append(propagators, propagation.TraceContext{})
		}
	}
	if len(propagators) == 0 {
		return propagation.TraceContext{}
	}
	return propagation.NewCompositeTextMapPropagator(propagators...)
}

// Init initializes opentelemetry tracing and returns a shutdown function to flush data immediately
// on a termination signal.
func Init(cfg *config.AgentConfig) func() {
	sdkconfig.InitConfig(cfg)

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !cfg.GetReporting().GetSecure().GetValue()},
	}}

	zipkinBatchExporter, err := zipkin.NewRawExporter(
		cfg.GetReporting().GetEndpoint().GetValue(),
		cfg.GetServiceName().GetValue(),
		zipkin.WithClient(client),
	)
	if err != nil {
		log.Fatal(err)
	}

	resources, err := resource.New(context.Background(), resource.WithAttributes(
		semconv.ServiceNameKey.String(cfg.GetServiceName().GetValue()),
		semconv.TelemetrySDKNameKey.String("hypertrace"),
		semconv.TelemetrySDKVersionKey.String(version.Version),
	))
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(zipkinBatchExporter, sdktrace.WithBatchTimeout(batchTimeoutInSecs*time.Millisecond)),
		sdktrace.WithResource(resources),
	)
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(makePropagator(cfg.PropagationFormats))

	// TODO: use batcher instead of this hack:
	return func() {
		// This is a sad temporary solution for the lack of flusher in the batcher interface.
		// What we do here is that we wait for `batchTimeout` seconds as that is the time configured
		// in the batcher and hence we make sure spans had time to be flushed.
		// In next versions the flush functionality is finally added and we will use it.
		<-time.After(batchTimeoutInSecs * 1.5 * time.Millisecond)
	}
}
