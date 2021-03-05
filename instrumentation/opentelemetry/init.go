package opentelemetry

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/hypertrace/goagent/sdk"

	"go.opentelemetry.io/otel/attribute"

	"crypto/tls"

	"github.com/hypertrace/goagent/config"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	"github.com/hypertrace/goagent/version"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

var batchTimeout = time.Duration(200) * time.Millisecond

var (
	traceProviders    map[string]*sdktrace.TracerProvider
	globalSampler     sdktrace.Sampler
	initialized       = false
	mu                sync.Mutex
	reportingEndpoint string
	secure            bool
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
	reportingEndpoint = cfg.GetReporting().GetEndpoint().GetValue()
	secure = cfg.GetReporting().GetSecure().GetValue()

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !secure},
	}}

	zipkinExporter, err := zipkin.NewRawExporter(
		reportingEndpoint,
		cfg.GetServiceName().GetValue(),
		zipkin.WithClient(client),
	)

	if cfg.Reporting.TraceReporterType == config.TraceReporterType_OTLP {
		opts := []otlpgrpc.Option{
			otlpgrpc.WithEndpoint(cfg.GetReporting().GetEndpoint().GetValue()),
		}

		if !cfg.GetReporting().GetSecure().GetValue() {
			opts = append(opts, otlpgrpc.WithInsecure())
		}

		batcherExporter, err = otlp.NewExporter(
			context.Background(),
			otlpgrpc.NewDriver(opts...),
		)
	} else {
		client := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !cfg.GetReporting().GetSecure().GetValue()},
		}}

		batcherExporter, err = zipkin.NewRawExporter(
			cfg.GetReporting().GetEndpoint().GetValue(),
			cfg.GetServiceName().GetValue(),
			zipkin.WithClient(client),
		)
	}
	if err != nil {
		log.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(zipkinExporter, sdktrace.WithBatchTimeout(batchTimeout))

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(createResources(cfg.GetServiceName().GetValue(), cfg.ResourceAttributes)...),
	)
	if err != nil {
		log.Fatal(err)
	}
	sampler := sdktrace.AlwaysSample()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}),
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(makePropagator(cfg.PropagationFormats))

	traceProviders = make(map[string]*sdktrace.TracerProvider)
	globalSampler = sampler
	initialized = true
	return func() {
		mu.Lock()
		defer mu.Unlock()
		for serviceName, tracerProvider := range traceProviders {
			tracerProvider.Shutdown(context.Background())
			delete(traceProviders, serviceName)
		}
		tp.Shutdown(context.Background())
		initialized = false
	}
}

func createResources(serviceName string, resources map[string]string) []attribute.KeyValue {
	retValues := []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
		semconv.TelemetrySDKNameKey.String("hypertrace"),
		semconv.TelemetrySDKVersionKey.String(version.Version),
	}

	for k, v := range resources {
		retValues = append(retValues, attribute.String(k, v))
	}
	return retValues
}

// RegisterService creates tracerprovider for a new service and returns a func which can be used to create spans
func RegisterService(serviceName string, resourceAttributes map[string]string) (sdk.StartSpan, error) {
	mu.Lock()
	defer mu.Unlock()
	if !initialized {
		return nil, fmt.Errorf("hypertrace lib not initialized. hypertrace.Init has not been called")
	}

	if _, ok := traceProviders[serviceName]; ok {
		return nil, fmt.Errorf("service %v already initialized", serviceName)
	}

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !secure},
	}}

	zipkinExporter, err := zipkin.NewRawExporter(
		reportingEndpoint,
		serviceName,
		zipkin.WithClient(client),
	)
	if err != nil {
		log.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(zipkinExporter)

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(createResources(serviceName, resourceAttributes)...),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: globalSampler}),
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resources),
	)

	traceProviders[serviceName] = tp
	return startSpan(func() trace.TracerProvider {
		return tp
	}), nil
}
