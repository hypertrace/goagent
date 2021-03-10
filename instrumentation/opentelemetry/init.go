package opentelemetry

import (
	"context"
	"fmt"
	"github.com/hypertrace/goagent/sdk"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel/label"

	"crypto/tls"
	"github.com/hypertrace/goagent/config"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	"github.com/hypertrace/goagent/version"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/propagation"
	sdkexport "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

var batchTimeout = time.Duration(200) * time.Millisecond

type traceProviderWrapper struct {
	tp       trace.TracerProvider
	shutdown func()
}

var (
	traceProviders map[string]*traceProviderWrapper
	batcher        sdkexport.SpanExporter
	sampler        sdktrace.Sampler
	initialized    = false
	mu             sync.Mutex
)

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
	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return func() {}
	}
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

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(createResources(cfg.GetServiceName().GetValue(), cfg.ResourceAttributes)...),
	)
	if err != nil {
		log.Fatal(err)
	}
	defaultSampler := sdktrace.AlwaysSample()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: defaultSampler}),
		sdktrace.WithBatcher(zipkinBatchExporter, sdktrace.WithBatchTimeout(batchTimeout)),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(makePropagator(cfg.PropagationFormats))

	batcher = zipkinBatchExporter
	traceProviders = make(map[string]*traceProviderWrapper)
	sampler = defaultSampler
	initialized = true
	return func() {
		mu.Lock()
		defer mu.Unlock()
		// TODO: There is an issue here. The traceprovider which is shutdown first, will only have its span flushed.
		for _, wrapper := range traceProviders {
			wrapper.shutdown()
		}
		tp.Shutdown(context.Background())
		initialized = false
	}
}

func createResources(serviceName string, resources map[string]string) []label.KeyValue {
	retValues := []label.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
		semconv.TelemetrySDKNameKey.String("hypertrace"),
		semconv.TelemetrySDKVersionKey.String(version.Version),
	}

	for k, v := range resources {
		retValues = append(retValues, label.String(k, v))
	}
	return retValues
}

func InitService(serviceName string, resourceAttributes map[string]string) (sdk.StartSpan, error) {
	mu.Lock()
	defer mu.Unlock()
	if !initialized {
		return nil, fmt.Errorf("hypertrace lib not initialized. hypertrace.Init has not been called")
	}

	if _, ok := traceProviders[serviceName]; ok {
		return nil, fmt.Errorf("service %v already initialized", serviceName)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(createResources(serviceName, resourceAttributes)...),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}),
		sdktrace.WithBatcher(batcher, sdktrace.WithBatchTimeout(batchTimeout)),
		sdktrace.WithResource(resources),
	)

	traceProviders[serviceName] = &traceProviderWrapper{
		tp: tp,
		shutdown: func() {
			tp.Shutdown(context.Background())
		},
	}
	return startSpan(tp), nil
}
