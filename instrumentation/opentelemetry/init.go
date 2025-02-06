package opentelemetry // import "github.com/hypertrace/goagent/instrumentation/opentelemetry"

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/zapr"
	config "github.com/hypertrace/agent-config/gen/go/v1"
	modbsp "github.com/hypertrace/goagent/instrumentation/opentelemetry/batchspanprocessor"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/errorhandler"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/identifier"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal/metrics"
	"github.com/hypertrace/goagent/sdk"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	"github.com/hypertrace/goagent/version"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otlphttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
)

var (
	batchTimeout          = time.Duration(200) * time.Millisecond
	traceProviders        map[string]*sdktrace.TracerProvider
	globalSampler         sdktrace.Sampler
	initialized           = false
	enabled               = false
	mu                    sync.Mutex
	exporterFactory       func() (sdktrace.SpanExporter, error)
	configFactory         func() *config.AgentConfig
	versionInfoAttributes = []attribute.KeyValue{
		semconv.TelemetrySDKNameKey.String("hypertrace"),
		semconv.TelemetrySDKVersionKey.String(version.Version),
	}
)

const AgentTokenKey = "traceableai-agent-token"

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

func makeMetricsExporterFactory(cfg *config.AgentConfig) func() (metric.Exporter, error) {
	// We are only supporting logging and otlp metric exporters for now. We will add support for prometheus
	// metrics later
	switch cfg.Reporting.MetricReporterType {
	case config.MetricReporterType_METRIC_REPORTER_TYPE_LOGGING:
		// stdout exporter
		return func() (metric.Exporter, error) {
			// TODO: Define if endpoint could be a filepath to write into a file.
			return stdoutmetric.New()
		}
	default:
		endpoint := cfg.GetReporting().GetMetricEndpoint().GetValue()
		if len(endpoint) == 0 {
			endpoint = cfg.GetReporting().GetEndpoint().GetValue()
		}

		opts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(removeProtocolPrefixForOTLP(endpoint)),
		}

		if !cfg.GetReporting().GetSecure().GetValue() {
			opts = append(opts, otlpmetricgrpc.WithInsecure())
		}

		certFile := cfg.GetReporting().GetCertFile().GetValue()
		if len(certFile) > 0 {
			if tlsCredentials, err := credentials.NewClientTLSFromFile(certFile, ""); err == nil {
				opts = append(opts, otlpmetricgrpc.WithTLSCredentials(tlsCredentials))
			} else {
				log.Printf("error while creating tls credentials from cert path %s: %v", certFile, err)
			}
		}

		if cfg.Reporting.GetEnableGrpcLoadbalancing().GetValue() {
			resolver.SetDefaultScheme("dns")
			opts = append(opts, otlpmetricgrpc.WithServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`))
		}

		return func() (metric.Exporter, error) {
			return otlpmetricgrpc.New(context.Background(), opts...)
		}
	}
}

func makeExporterFactory(cfg *config.AgentConfig) func() (sdktrace.SpanExporter, error) {
	additionalHeaders := make(map[string]string)
	if cfg.Reporting.GetToken() != nil {
		additionalHeaders[AgentTokenKey] = cfg.GetReporting().GetToken().GetValue()
	}

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
				zipkin.WithHeaders(additionalHeaders),
			)
		}
	case config.TraceReporterType_LOGGING:
		return func() (sdktrace.SpanExporter, error) {
			// TODO: Define if endpoint could be a filepath to write into a file.
			return stdouttrace.New(stdouttrace.WithPrettyPrint())
		}

	case config.TraceReporterType_OTLP_HTTP:
		opts := []otlphttp.Option{
			otlphttp.WithEndpoint(cfg.GetReporting().GetEndpoint().GetValue()),
		}

		if !cfg.GetReporting().GetSecure().GetValue() {
			opts = append(opts, otlphttp.WithInsecure())
		}
		opts = append(opts, otlphttp.WithHeaders(additionalHeaders))

		certFile := cfg.GetReporting().GetCertFile().GetValue()
		if len(certFile) > 0 {
			opts = append(opts, otlphttp.WithTLSClientConfig(createTLSConfig(cfg.GetReporting())))
		}

		return func() (sdktrace.SpanExporter, error) {
			return otlphttp.New(context.Background(), opts...)
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

		opts = append(opts, otlpgrpc.WithHeaders(additionalHeaders))

		return func() (sdktrace.SpanExporter, error) {
			return otlptrace.New(
				context.Background(),
				otlpgrpc.NewClient(opts...),
			)
		}
	}
}

func makeConfigFactory(cfg *config.AgentConfig) func() *config.AgentConfig {
	return func() *config.AgentConfig {
		return cfg
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
	return InitWithSpanProcessorWrapper(cfg, nil, versionInfoAttributes)
}

// InitWithSpanProcessorWrapper initializes opentelemetry tracing with a wrapper over span processor
// and returns a shutdown function to flush data immediately on a termination signal.
func InitWithSpanProcessorWrapper(cfg *config.AgentConfig, wrapper SpanProcessorWrapper,
	versionInfoAttrs []attribute.KeyValue) func() {
	logger, err := zap.NewProduction()
	if err != nil {
		logger = nil
		log.Printf("error while creating default zap logger %v", err)
	}
	return InitWithSpanProcessorWrapperAndZap(cfg, wrapper, versionInfoAttrs, logger)
}

// InitWithSpanProcessorWrapperAndZap initializes opentelemetry tracing with a wrapper over span processor
// and returns a shutdown function to flush data immediately on a termination signal.
// Also sets opentelemetry internal errorhandler to the provider zap errorhandler
func InitWithSpanProcessorWrapperAndZap(cfg *config.AgentConfig, wrapper SpanProcessorWrapper,
	versionInfoAttrs []attribute.KeyValue, logger *zap.Logger) func() {
	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return func() {}
	}
	sdkconfig.InitConfig(cfg)

	enabled = cfg.GetEnabled().Value
	if !enabled {
		initialized = true
		otel.SetTracerProvider(noop.NewTracerProvider())
		// even if the tracer isn't enabled, propagation is still enabled
		// to not break the full workflow of the tracing system. Even
		// if this service will not report spans and the trace might look
		// broken, spans can still be grouped by trace ID.
		otel.SetTextMapPropagator(makePropagator(cfg.PropagationFormats))
		return func() {
			initialized = false
			sdkconfig.ResetConfig()
		}
	}

	if logger != nil {
		_ = zap.ReplaceGlobals(logger.With(zap.String("service", "hypertrace")))

		// initialize opentelemetry's internal logger
		logr := zapr.NewLogger(logger)
		otel.SetLogger(logr)

		// initialize opentelemetry's internal error handler
		errorhandler.Init(logger)
	}

	// Initialize metrics
	metricsShutdownFn := initializeMetrics(cfg, versionInfoAttrs)

	exporterFactory = makeExporterFactory(cfg)
	configFactory = makeConfigFactory(cfg)

	exporter, err := exporterFactory()
	if err != nil {
		log.Fatal(err)
	}

	sp := modbsp.CreateBatchSpanProcessor(
		shouldUseCustomBatchSpanProcessor(cfg),
		exporter,
		sdktrace.WithBatchTimeout(batchTimeout))
	if wrapper != nil {
		sp = &spanProcessorWithWrapper{wrapper, sp}
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(createResources(getResourceAttrsWithServiceName(cfg.ResourceAttributes, cfg.GetServiceName().GetValue()),
			versionInfoAttrs)...),
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
		for key, tracerProvider := range traceProviders {
			err := tracerProvider.Shutdown(context.Background())
			if err != nil {
				log.Printf("error while shutting down tracer provider: %v\n", err)
			}
			delete(traceProviders, key)
		}
		traceProviders = map[string]*sdktrace.TracerProvider{}
		err := tp.Shutdown(context.Background())
		if err != nil {
			log.Printf("error while shutting down default tracer provider: %v\n", err)
		}

		metricsShutdownFn()
		initialized = false
		enabled = false
		sdkconfig.ResetConfig()
	}
}

func createResources(resources map[string]string,
	versionInfo []attribute.KeyValue) []attribute.KeyValue {
	retValues := []attribute.KeyValue{
		semconv.TelemetrySDKLanguageGo,
	}

	retValues = append(retValues, versionInfo...)

	for k, v := range resources {
		retValues = append(retValues, attribute.String(k, v))
	}
	return retValues
}

// RegisterService creates tracerprovider for a new service (represented via a unique key) and returns a func which can be used to create spans and the TracerProvider
func RegisterService(key string, resourceAttributes map[string]string) (sdk.StartSpan, trace.TracerProvider, error) {
	return RegisterServiceWithSpanProcessorWrapper(key, resourceAttributes, nil, versionInfoAttributes)
}

// RegisterServiceWithSpanProcessorWrapper creates a tracerprovider for a new service (represented via a unique key) with a wrapper over opentelemetry span processor
// and returns a func which can be used to create spans and the TracerProvider
func RegisterServiceWithSpanProcessorWrapper(key string, resourceAttributes map[string]string,
	wrapper SpanProcessorWrapper, versionInfoAttrs []attribute.KeyValue) (sdk.StartSpan, trace.TracerProvider, error) {
	mu.Lock()
	defer mu.Unlock()
	if !initialized {
		return nil, noop.NewTracerProvider(), fmt.Errorf("hypertrace hadn't been initialized")
	}

	if !enabled {
		return NoopStartSpan, noop.NewTracerProvider(), nil
	}

	if _, ok := traceProviders[key]; ok {
		return nil, noop.NewTracerProvider(), fmt.Errorf("key %v is already used for initialization", key)
	}

	exporter, err := exporterFactory()
	if err != nil {
		log.Fatal(err)
	}

	sp := modbsp.CreateBatchSpanProcessor(
		shouldUseCustomBatchSpanProcessor(configFactory()),
		exporter,
		sdktrace.WithBatchTimeout(batchTimeout))
	if wrapper != nil {
		sp = &spanProcessorWithWrapper{wrapper, sp}
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(createResources(resourceAttributes, versionInfoAttrs)...),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(globalSampler),
		sdktrace.WithSpanProcessor(sp),
		sdktrace.WithResource(resources),
	)

	traceProviders[key] = tp
	return startSpan(func() trace.TracerProvider {
		return tp
	}), tp, nil
}

func initializeMetrics(cfg *config.AgentConfig, versionInfoAttrs []attribute.KeyValue) func() {
	if shouldDisableMetrics(cfg) {
		return func() {}
	}

	metricsExporterFactory := makeMetricsExporterFactory(cfg)
	metricsExporter, err := metricsExporterFactory()
	if err != nil {
		log.Fatal(err)
	}
	periodicReader := metric.NewPeriodicReader(metricsExporter)

	resourceKvps := createResources(getResourceAttrsWithServiceName(cfg.ResourceAttributes, cfg.GetServiceName().GetValue()), versionInfoAttrs)
	resourceKvps = append(resourceKvps, identifier.ServiceInstanceKeyValue)
	metricResources, err := resource.New(context.Background(), resource.WithAttributes(resourceKvps...))
	if err != nil {
		log.Fatal(err)
	}
	meterProvider := metric.NewMeterProvider(metric.WithReader(periodicReader), metric.WithResource(metricResources))
	otel.SetMeterProvider(meterProvider)

	metrics.InitializeSystemMetrics()
	return func() {
		err = meterProvider.Shutdown(context.Background())
		if err != nil {
			log.Printf("an error while calling metrics provider shutdown: %v", err)
		}
		err := periodicReader.Shutdown(context.Background())
		if err != nil {
			log.Printf("an error while calling metrics reader shutdown: %v", err)
		}
	}
}

func shouldDisableMetrics(cfg *config.AgentConfig) bool {
	// Disable metrics if the tracing exporter is not OTLP(grpc) and the metrics endpoint is not explicitly set.
	// This is because we use the traces OTLP endpoint for metrics if the metrics endpoint is not set.
	// By default the traces endpoint is zipkin which does not have support for metrics.
	if cfg.GetReporting() != nil && cfg.GetReporting().GetTraceReporterType() != config.TraceReporterType_OTLP &&
		len(cfg.GetReporting().GetMetricEndpoint().GetValue()) == 0 {
		return true
	}

	return cfg.GetTelemetry() == nil || !cfg.GetTelemetry().GetMetricsEnabled().GetValue()
}

func shouldUseCustomBatchSpanProcessor(cfg *config.AgentConfig) bool {
	return (cfg.GetGoagent() != nil && cfg.GetGoagent().GetUseCustomBsp().GetValue()) && // bsp enabled AND
		(cfg.GetTelemetry() != nil && cfg.GetTelemetry().GetMetricsEnabled().GetValue()) // metrics enabled
}

func getResourceAttrsWithServiceName(resourceMap map[string]string, serviceName string) map[string]string {
	if resourceMap == nil {
		resourceMap = make(map[string]string)
	}
	serviceNameKey := string(semconv.ServiceNameKey)
	if _, ok := resourceMap[serviceNameKey]; !ok && (len(serviceName) > 0) {
		resourceMap[serviceNameKey] = serviceName
	}

	return resourceMap
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
