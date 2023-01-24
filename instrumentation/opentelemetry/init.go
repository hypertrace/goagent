package opentelemetry // import "github.com/hypertrace/goagent/instrumentation/opentelemetry"

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/trace"

	"github.com/hypertrace/goagent/sdk"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"go.opentelemetry.io/otel/attribute"

	"github.com/go-logr/stdr"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	"github.com/hypertrace/goagent/version"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc/credentials"
	//"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	//"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	//"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	// otelmetric "go.opentelemetry.io/otel/metric"
	otelmetricglobal "go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	// "go.opentelemetry.io/otel/sdk/metric"
	// "go.opentelemetry.io/otel/sdk/metric/metricdata"
	// "go.opentelemetry.io/otel/sdk/metric/view"
)

var batchTimeout = time.Duration(200) * time.Millisecond

var (
	traceProviders  map[string]*sdktrace.TracerProvider
	globalSampler   sdktrace.Sampler
	initialized     = false
	enabled         = false
	mu              sync.Mutex
	exporterFactory func() (sdktrace.SpanExporter, error)
)

func makePropagator(formats []config.PropagationFormat) propagation.TextMapPropagator {
	var propagators []propagation.TextMapPropagator
	for _, format := range formats {
		switch format {
		case config.PropagationFormat_B3:
			// We set B3MultipleHeader in here but ideally we should use both.
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader|b3.B3SingleHeader)))
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

func makeExporterFactory(cfg *config.AgentConfig) func() (sdktrace.SpanExporter, error) {
	switch cfg.Reporting.TraceReporterType {
	case config.TraceReporterType_ZIPKIN:
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: createTLSConfig(cfg.GetReporting()),
			},
		}

		return func() (sdktrace.SpanExporter, error) {
			return zipkin.New(
				cfg.GetReporting().GetEndpoint().GetValue(),
				zipkin.WithClient(client),
			)
		}
	case config.TraceReporterType_LOGGING:
		return func() (sdktrace.SpanExporter, error) {
			// TODO: Define if endpoint could be a filepath to write into a file.
			return stdouttrace.New(stdouttrace.WithPrettyPrint())
		}
	default:
		opts := []otlpgrpc.Option{
			otlpgrpc.WithEndpoint(removeProtocolPrefixForOTLP(cfg.GetReporting().GetEndpoint().GetValue())),
		}

		if !cfg.GetReporting().GetSecure().GetValue() {
			opts = append(opts, otlpgrpc.WithInsecure())
		}

		certFile := cfg.GetReporting().GetCertFile().GetValue()
		if len(certFile) > 0 {
			if tlsCredentials, err := credentials.NewClientTLSFromFile(certFile, ""); err == nil {
				opts = append(opts, otlpgrpc.WithTLSCredentials(tlsCredentials))
			} else {
				log.Printf("error while creating tls credentials from cert path %s: %v", certFile, err)
			}
		}

		if cfg.Reporting.GetEnableGrpcLoadbalancing().GetValue() {
			resolver.SetDefaultScheme("dns")
			opts = append(opts, otlpgrpc.WithServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`))
		}

		return func() (sdktrace.SpanExporter, error) {
			return otlptrace.New(
				context.Background(),
				otlpgrpc.NewClient(opts...),
			)
		}
	}
}

func createTLSConfig(reportingCfg *config.Reporting) *tls.Config {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	tlsConfig.InsecureSkipVerify = !reportingCfg.GetSecure().GetValue()
	certFile := reportingCfg.GetCertFile().GetValue()
	if len(certFile) > 0 {
		tlsConfig.RootCAs = createCaCertPoolFromFile(certFile)
	}

	return tlsConfig
}

// createCaCertPoolFromFile creates a CA Cert Pool from a file path containing
// a raw CA certificate to verify a server certificate. The file path is the
// reporting.cert_file config value.
func createCaCertPoolFromFile(certFile string) *x509.CertPool {
	certBytes, err := os.ReadFile(filepath.Clean(certFile))
	if err != nil {
		log.Printf("error while reading cert path: %v", err)
		return nil
	}
	cp := x509.NewCertPool()
	if ok := cp.AppendCertsFromPEM(certBytes); !ok {
		log.Printf("error while configuring tls: failed to append certificate to the cert pool")
		return nil
	}

	return cp
}

// Init initializes opentelemetry tracing and returns a shutdown function to flush data immediately
// on a termination signal.
func Init(cfg *config.AgentConfig) func() {
	return InitWithSpanProcessorWrapper(cfg, nil)
}

// InitWithSpanProcessorWrapper initializes opentelemetry tracing with a wrapper over span processor
// and returns a shutdown function to flush data immediately on a termination signal.
func InitWithSpanProcessorWrapper(cfg *config.AgentConfig, wrapper SpanProcessorWrapper) func() {
	stdr.SetVerbosity(5)
	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return func() {}
	}
	sdkconfig.InitConfig(cfg)

	enabled = cfg.GetEnabled().Value
	if !enabled {
		initialized = true
		otel.SetTracerProvider(trace.NewNoopTracerProvider())
		// even if the tracer isn't enabled, propagation is still enabled
		// to not to break the full workflow of the tracing system. Even
		// if this service will not report spans and the trace might look
		// broken, spans can still be grouped by trace ID.
		otel.SetTextMapPropagator(makePropagator(cfg.PropagationFormats))
		return func() {
			initialized = false
			sdkconfig.ResetConfig()
		}
	}

	exporterFactory = makeExporterFactory(cfg)

	exporter, err := exporterFactory()
	if err != nil {
		log.Fatal(err)
	}

	sp := sdktrace.NewBatchSpanProcessor(exporter, sdktrace.WithBatchTimeout(batchTimeout))
	if wrapper != nil {
		sp = &spanProcessorWithWrapper{wrapper, sp}
	}

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
		sdktrace.WithSpanProcessor(sp),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(makePropagator(cfg.PropagationFormats))

	initMetrics()

	traceProviders = make(map[string]*sdktrace.TracerProvider)
	globalSampler = sampler
	initialized = true

	startSpanFn := startSpan(func() trace.TracerProvider {
		return tp
	})

	// Startup span
	if cfg.GetTelemetry().GetStartupSpanEnabled().GetValue() {
		_, span, ender := startSpanFn(context.Background(), "startup", &sdk.SpanOptions{})
		span.SetAttribute("hypertrace.agent.startup", true)
		ender()
	}

	return func() {
		mu.Lock()
		defer mu.Unlock()
		for serviceName, tracerProvider := range traceProviders {
			err := tracerProvider.Shutdown(context.Background())
			if err != nil {
				log.Printf("error while shutting down tracer provider: %v\n", err)
			}
			delete(traceProviders, serviceName)
		}
		traceProviders = map[string]*sdktrace.TracerProvider{}
		err := tp.Shutdown(context.Background())
		if err != nil {
			log.Printf("error while shutting down default tracer provider: %v\n", err)
		}
		initialized = false
		enabled = false
		sdkconfig.ResetConfig()
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

// RegisterService creates tracerprovider for a new service and returns a func which can be used to create spans and the TracerProvider
func RegisterService(serviceName string, resourceAttributes map[string]string) (sdk.StartSpan, trace.TracerProvider, error) {
	return RegisterServiceWithSpanProcessorWrapper(serviceName, resourceAttributes, nil)
}

// RegisterServiceWithSpanProcessorWrapper creates a tracerprovider for a new service with a wrapper over opentelemetry span processor
// and returns a func which can be used to create spans and the TracerProvider
func RegisterServiceWithSpanProcessorWrapper(serviceName string, resourceAttributes map[string]string,
	wrapper SpanProcessorWrapper) (sdk.StartSpan, trace.TracerProvider, error) {
	mu.Lock()
	defer mu.Unlock()
	if !initialized {
		return nil, trace.NewNoopTracerProvider(), fmt.Errorf("hypertrace hadn't been initialized")
	}

	if !enabled {
		return NoopStartSpan, trace.NewNoopTracerProvider(), nil
	}

	if _, ok := traceProviders[serviceName]; ok {
		return nil, trace.NewNoopTracerProvider(), fmt.Errorf("service %v already initialized", serviceName)
	}

	exporter, err := exporterFactory()
	if err != nil {
		log.Fatal(err)
	}

	sp := sdktrace.NewBatchSpanProcessor(exporter, sdktrace.WithBatchTimeout(batchTimeout))
	if wrapper != nil {
		sp = &spanProcessorWithWrapper{wrapper, sp}
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(createResources(serviceName, resourceAttributes)...),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(globalSampler),
		sdktrace.WithSpanProcessor(sp),
		sdktrace.WithResource(resources),
	)

	traceProviders[serviceName] = tp
	return startSpan(func() trace.TracerProvider {
		return tp
	}), tp, nil
}

// SpanProcessorWrapper wraps otel span processor
// and is responsible to delegate calls to the wrapped processor
type SpanProcessorWrapper interface {
	OnStart(parent context.Context, s sdktrace.ReadWriteSpan, delegate sdktrace.SpanProcessor)
	OnEnd(s sdktrace.ReadOnlySpan, delegate sdktrace.SpanProcessor)
}

type spanProcessorWithWrapper struct {
	wrapper   SpanProcessorWrapper
	processor sdktrace.SpanProcessor
}

func (sp *spanProcessorWithWrapper) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {
	sp.wrapper.OnStart(parent, s, sp.processor)
}

func (sp *spanProcessorWithWrapper) OnEnd(s sdktrace.ReadOnlySpan) {
	sp.wrapper.OnEnd(s, sp.processor)
}

func (sp *spanProcessorWithWrapper) Shutdown(ctx context.Context) error {
	return sp.processor.Shutdown(ctx)
}

func (sp *spanProcessorWithWrapper) ForceFlush(ctx context.Context) error {
	return sp.processor.ForceFlush(ctx)
}

func initMetrics() {
	// stdout exporter
	// exporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	// if err != nil {
	// 	log.Printf("error in init metrics: %v", fmt.Errorf("creating stdoutmetric exporter: %w", err))
	// 	//return nil, fmt.Errorf("creating stdoutmetric exporter: %w", err)
	// return
	// }

	// otlp exporter
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	}

	exporter, err := otlpmetric.New(
		context.Background(),
		otlpmetricgrpc.NewClient(opts...),
	)
	if err != nil {
		log.Printf("error in init metrics: %v", fmt.Errorf("creating otlpmetric exporter: %w", err))
		return
	}

	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithInexpensiveDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
	)
	if err := pusher.Start(context.Background()); err != nil {
		log.Fatalf("starting push controller: %v", err)
	}

	otelmetricglobal.SetMeterProvider(pusher)

	// metricsClient :=
	// defaultView, _ := view.New(view.MatchInstrumentName("*"))

	// meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(metricsClient,
	// 	metric.WithAggregationSelector(metric.DefaultAggregationSelector),
	// 	metric.WithTemporalitySelector(deltaTemporalitySelector),
	// ), defaultView, defaultView))

	// otelmetricglobal.SetMeterProvider(meterProvider)
}
