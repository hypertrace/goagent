package opentelemetry

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/hypertrace/goagent/version"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/semconv"

	"crypto/tls"

	"github.com/hypertrace/goagent/config"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var batchTimeout = time.Duration(200) * time.Millisecond

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
	return InitWithResources(cfg, nil)
}

// Init initializes opentelemetry tracing and returns a shutdown function to flush data immediately
// on a termination signal.
func InitWithResources(cfg *config.AgentConfig, initResources map[string]interface{}) func() {
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
		getResources(initResources, semconv.ServiceNameKey.String(cfg.GetServiceName().GetValue()),
			semconv.TelemetrySDKNameKey.String("hypertrace"),
			semconv.TelemetrySDKVersionKey.String(version.Version))...,
	))
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(zipkinBatchExporter, sdktrace.WithBatchTimeout(batchTimeout)),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(makePropagator(cfg.PropagationFormats))

	return func() {
		tp.Shutdown(context.Background())
	}
}

func getResources(initResources map[string]interface{}, values ...label.KeyValue) []label.KeyValue {
	irl := len(initResources)
	if irl == 0 {
		return values
	}
	retValues := make([]label.KeyValue, irl+len(values))
	for _, value := range values {
		retValues = append(retValues, value)
	}
	for k, v := range initResources {
		retValues = append(retValues, label.Any(k, v))
	}

	return retValues
}
