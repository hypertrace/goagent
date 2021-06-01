package opentelemetry

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

var batchTimeout = time.Duration(200) * time.Millisecond

var (
	traceProviders  map[string]*sdktrace.TracerProvider
	globalSampler   sdktrace.Sampler
	initialized     = false
	mu              sync.Mutex
	exporterFactory func(serviceName string) (export.SpanExporter, error)
)

func makePropagator(formats []config.PropagationFormat) propagation.TextMapPropagator {
	var propagators []propagation.TextMapPropagator
	for _, format := range formats {
		switch format {
		case config.PropagationFormat_B3:
			// We set B3MultipleHeader in here but ideally we should use both.
			propagators = append(propagators, b3.B3{InjectEncoding: b3.B3MultipleHeader | b3.B3SingleHeader})
		case config.PropagationFormat_TRACECONTEXT:
			propagators = append(propagators, propagation.TraceContext{})
		}
	}
	if len(propagators) == 0 {
		return propagation.TraceContext{}
	}
	return propagation.NewCompositeTextMapPropagator(propagators...)
}

// removeProtocolPrefixForOTLP removes the prefix protocol as grpc exporter
// will reject it with error "too many colons in address"
func removeProtocolPrefixForOTLP(endpoint string) string {
	pieces := strings.SplitN(endpoint, "://", 2)
	if len(pieces) == 1 {
		return endpoint
	}

	return pieces[1]
}

func makeExporterFactory(cfg *config.AgentConfig) func(serviceName string) (export.SpanExporter, error) {
	switch cfg.Reporting.TraceReporterType {
	case config.TraceReporterType_ZIPKIN:
		client := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !cfg.GetReporting().GetSecure().GetValue()},
		}}

		return func(serviceName string) (export.SpanExporter, error) {
			return zipkin.NewRawExporter(
				cfg.GetReporting().GetEndpoint().GetValue(),
				zipkin.WithClient(client),
			)
		}
	default:
		opts := []otlpgrpc.Option{
			otlpgrpc.WithEndpoint(removeProtocolPrefixForOTLP(cfg.GetReporting().GetEndpoint().GetValue())),
		}

		if !cfg.GetReporting().GetSecure().GetValue() {
			opts = append(opts, otlpgrpc.WithInsecure())
		}

		return func(_ string) (export.SpanExporter, error) {
			return otlp.NewExporter(
				context.Background(),
				otlpgrpc.NewDriver(opts...),
			)
		}
	}
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

	exporterFactory = makeExporterFactory(cfg)

	exporter, err := exporterFactory(cfg.ServiceName.Value)
	if err != nil {
		log.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter, sdktrace.WithBatchTimeout(batchTimeout))

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(createResources(cfg.GetServiceName().GetValue(), cfg.ResourceAttributes)...),
	)
	if err != nil {
		log.Fatal(err)
	}

	sampler := sdktrace.AlwaysSample()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
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
		semconv.TelemetrySDKLanguageGo,
	}

	for k, v := range resources {
		retValues = append(retValues, attribute.String(k, v))
	}
	return retValues
}

// RegisterService creates tracerprovider for a new service and returns a func which can be used to create spans
func RegisterService(serviceName string, resourceAttributes map[string]string) (
	sdk.StartSpan,
	sdk.StartReadbackSpan,
	error,
) {
	mu.Lock()
	defer mu.Unlock()
	if !initialized {
		return nil, nil, fmt.Errorf("hypertrace hadn't been initialized")
	}

	if _, ok := traceProviders[serviceName]; ok {
		return nil, nil, fmt.Errorf("service %v already initialized", serviceName)
	}

	exporter, err := exporterFactory(serviceName)
	if err != nil {
		log.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter, sdktrace.WithBatchTimeout(batchTimeout))

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(createResources(serviceName, resourceAttributes)...),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(globalSampler),
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resources),
	)

	traceProviders[serviceName] = tp
	tper := func() trace.TracerProvider {
		return tp
	}

	return startSpan(tper), startReadbackSpan(tper), nil
}
