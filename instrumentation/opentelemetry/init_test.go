package opentelemetry

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/resolver"
)

func ExampleInit() {
	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.DataCapture.HttpHeaders.Request = config.Bool(true)
	cfg.Reporting.Endpoint = config.String("http://api.traceable.ai:9411/api/v2/spans")

	shutdown := Init(cfg)
	defer shutdown()
}

func ExampleRegisterService() {
	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.DataCapture.HttpHeaders.Request = config.Bool(true)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_LOGGING

	shutdown := Init(cfg)
	defer shutdown()

	_, _, err := RegisterService("custom_env", "custom_service", map[string]string{"test1": "val1"})
	if err != nil {
		log.Fatalf("Error while initializing service: %v", err)
	}
}

func TestMultiRegisterService(t *testing.T) {
	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.DataCapture.HttpHeaders.Request = config.Bool(true)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_LOGGING

	shutdown := Init(cfg)
	defer shutdown()

	_, _, err := RegisterService("custom_env", "custom_service", map[string]string{"test1": "val1"})
	assert.Nil(t, err)

	_, _, err = RegisterService("custom_env", "custom_service_2", map[string]string{"test2": "val2"})
	assert.Nil(t, err)

	_, _, err = RegisterService("custom_env_1", "custom_service", map[string]string{"test3": "val3"})
	assert.Nil(t, err)
}

func TestInitDisabledAgent(t *testing.T) {
	cfg := config.Load()
	cfg.Enabled = config.Bool(false)
	shutdown := Init(cfg)
	defer shutdown()

	startSpan, tp, err := RegisterService("custom_env", "test_service", nil)
	require.NoError(t, err)
	assert.Equal(t, noop.NewTracerProvider(), tp)
	_, s, _ := startSpan(context.Background(), "test_span", nil)
	require.NoError(t, err)
	assert.True(t, s.IsNoop())
}

func TestInitWithCertfileAndSecure(t *testing.T) {
	cfg := config.Load()
	cfg.Reporting.Secure = config.Bool(true)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP
	cfg.Reporting.CertFile = config.String("testdata/rootCA.crt")
	cfg.Enabled = config.Bool(true)

	shutdown := Init(cfg)
	defer shutdown()
}

func TestOtlpService(t *testing.T) {
	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.DataCapture.HttpHeaders.Request = config.Bool(true)
	cfg.Reporting.Endpoint = config.String("http://api.traceable.ai:4317")
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP
	cfg.Enabled = config.Bool(true)

	shutdown := Init(cfg)
	defer shutdown()

	startSpan, tp, err := RegisterService("custom_env", "custom_service", map[string]string{"test1": "val1"})
	_, s, _ := startSpan(context.Background(), "test_span", nil)
	assert.False(t, s.IsNoop())
	assert.NotEqual(t, noop.NewTracerProvider(), tp)
	assert.Len(t, s.GetAttributes().GetValue("service.instance.id"), 36)
	if err != nil {
		log.Fatalf("Error while initializing service: %v", err)
	}
}

func TestGrpcLoadBalancingConfig(t *testing.T) {
	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.Reporting.Endpoint = config.String("http://api.traceable.ai:4317")
	cfg.Reporting.EnableGrpcLoadbalancing = config.Bool(true)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP
	cfg.Enabled = config.Bool(true)

	shutdown := Init(cfg)
	defer shutdown()

	assert.Equal(t, resolver.GetDefaultScheme(), "dns")
	_, tp, err := RegisterService("custom_env", "custom_service", map[string]string{"test1": "val1"})
	assert.NotEqual(t, noop.NewTracerProvider(), tp)
	if err != nil {
		log.Fatalf("Error while initializing service: %v", err)
	}
}

func TestShutdownFlushesAllSpans(t *testing.T) {
	requestIsReceived := false
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		requestIsReceived = true
		rw.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.Reporting.Endpoint = config.String(srv.URL)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_ZIPKIN
	cfg.Enabled = config.Bool(true)

	// By doing this we make sure a batching isn't happening
	batchTimeout = time.Duration(10) * time.Second

	shutdown := Init(cfg)
	assert.True(t, initialized)
	assert.Equal(t, 0, len(traceProviders))

	_, _, spanEnder := StartSpan(context.Background(), "my_span", nil)
	spanEnder()

	assert.False(t, requestIsReceived)
	shutdown()
	assert.True(t, requestIsReceived)
}

func TestMultipleTraceProviders(t *testing.T) {
	count := 0
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		count++
		rw.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.Reporting.Endpoint = config.String(srv.URL)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_ZIPKIN
	cfg.Enabled = config.Bool(true)
	// Disable metrics to only test trace provider.
	cfg.Telemetry.MetricsEnabled = config.Bool(false)

	// By doing this we make sure a batching isn't happening
	batchTimeout = time.Duration(10) * time.Second

	shutdown := Init(cfg)

	assert.True(t, initialized)
	assert.Equal(t, 0, len(traceProviders))

	_, _, spanEnder := StartSpan(context.Background(), "example_span", nil)
	spanEnder()

	startServiceSpan, tp, err := RegisterService("custom_env", "custom_service", map[string]string{"test1": "val1"})
	assert.NoError(t, err)
	assert.NotNil(t, startServiceSpan)
	assert.True(t, initialized)
	assert.Equal(t, 1, len(traceProviders))
	assert.NotEqual(t, noop.NewTracerProvider(), tp)

	_, _, serviceSpanEnder := startServiceSpan(context.Background(), "my_span", nil)
	serviceSpanEnder()

	t.Run("test no requests before flush", func(t *testing.T) {
		assert.Equal(t, 0, count)
	})

	t.Run("test 2 requests after flush", func(t *testing.T) {
		shutdown()
		assert.Equal(t, 2, count)
		assert.Equal(t, 0, len(traceProviders))
	})
}

func TestMultipleTraceProvidersCallAfterShutdown(t *testing.T) {
	requestIsReceived := false
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		requestIsReceived = true
		rw.WriteHeader(http.StatusAccepted)

	}))
	defer srv.Close()

	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.Reporting.Endpoint = config.String(srv.URL)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_ZIPKIN
	cfg.Enabled = config.Bool(true)

	// By doing this we make sure a batching isn't happening
	batchTimeout = time.Duration(10) * time.Second

	shutdown := Init(cfg)
	assert.True(t, initialized)
	assert.Equal(t, 0, len(traceProviders))

	startServiceSpan, _, err := RegisterService("custom_env", "custom_service", map[string]string{"test1": "val1"})
	assert.NoError(t, err)
	assert.NotNil(t, startServiceSpan)
	assert.True(t, initialized)
	assert.Equal(t, 1, len(traceProviders))

	_, _, spanEnder := startServiceSpan(context.Background(), "my_span", nil)
	spanEnder()

	assert.False(t, requestIsReceived)
	shutdown()
	assert.True(t, requestIsReceived)

	_, _, spanEnder = startServiceSpan(context.Background(), "my_span1", nil)
	spanEnder()
}

type carrier struct {
	m map[string]string
}

func (c carrier) Get(string) string { return "" }

func (c carrier) Set(key string, value string) {
	c.m[key] = value
}

func (c carrier) Keys() []string {
	keys := make([]string, len(c.m))
	idx := 0
	for k := range c.m {
		keys[idx] = k
		idx++
	}
	return keys
}

var _ propagation.TextMapCarrier = carrier{}

func TestPropagationFormats(t *testing.T) {
	cfg := config.Load()
	cfg.PropagationFormats = config.PropagationFormats(
		config.PropagationFormat_B3,
		config.PropagationFormat_TRACECONTEXT,
	)
	cfg.Enabled = config.Bool(true)

	shutdown := Init(cfg)
	defer shutdown()

	tracer := otel.Tracer("b3")
	ctx, _ := tracer.Start(context.Background(), "test")
	propagator := otel.GetTextMapPropagator()
	c := carrier{make(map[string]string)}
	propagator.Inject(ctx, c)
	_, ok := c.m["x-b3-traceid"]
	assert.True(t, ok)
	_, ok = c.m["traceparent"]
	assert.True(t, ok)
}

func TestTraceReporterType(t *testing.T) {
	cfg := config.Load()
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP
	cfg.Enabled = config.Bool(true)

	shutdown := Init(cfg)
	defer shutdown()
}

func TestRemoveProtocolPrefixForOTLP(t *testing.T) {
	assert.Equal(
		t,
		"traceable-agent.traceableai:4317",
		removeProtocolPrefixForOTLP("http://traceable-agent.traceableai:4317"),
	)

	assert.Equal(
		t,
		"traceable-agent.traceableai:4317",
		removeProtocolPrefixForOTLP("traceable-agent.traceableai:4317"),
	)
}

func TestCreateTLSConfigNoCertFileAndInsecure(t *testing.T) {
	// Just using default values
	cfg := config.Load()
	tlsConfig := createTLSConfig(cfg.GetReporting())

	assert.True(t, tlsConfig.InsecureSkipVerify)
	assert.Nil(t, tlsConfig.RootCAs)
}

func TestCreateTLSConfigNoCertFileButSecure(t *testing.T) {
	cfg := config.Load()
	cfg.Reporting.Secure = config.Bool(true)
	tlsConfig := createTLSConfig(cfg.GetReporting())

	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.Nil(t, tlsConfig.RootCAs)
}

func TestCreateTLSConfigCertFilePresentButInsecure(t *testing.T) {
	cfg := config.Load()
	cfg.Reporting.CertFile = config.String("testdata/rootCA.crt")
	tlsConfig := createTLSConfig(cfg.GetReporting())

	assert.True(t, tlsConfig.InsecureSkipVerify)
	assert.NotNil(t, tlsConfig.RootCAs)
}

func TestCreateTLSConfigCertFilePresentAndSecure(t *testing.T) {
	cfg := config.Load()
	cfg.Reporting.Secure = config.Bool(true)
	cfg.Reporting.CertFile = config.String("testdata/rootCA.crt")
	tlsConfig := createTLSConfig(cfg.GetReporting())

	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.NotNil(t, tlsConfig.RootCAs)
}

func TestCreateCaCertPoolFromFileThatDoesNotExist(t *testing.T) {
	assert.Nil(t, createCaCertPoolFromFile("testdata/nonExistentCA.crt"))
}

func TestCreateCaCertPoolFromFileThatIsBogus(t *testing.T) {
	assert.Nil(t, createCaCertPoolFromFile("testdata/fakeRootCA.crt"))
}

type mockSpanProcessorWrapper struct {
	onStartCount int
	onEndCount   int
}

func (spw *mockSpanProcessorWrapper) OnStart(parent context.Context, s sdktrace.ReadWriteSpan, delegate sdktrace.SpanProcessor) {
	spw.onStartCount++
	delegate.OnStart(parent, s)
}

func (spw *mockSpanProcessorWrapper) OnEnd(s sdktrace.ReadOnlySpan, delegate sdktrace.SpanProcessor) {
	spw.onEndCount++
	delegate.OnEnd(s)
}

func TestInitWithSpanProcessorWrapper(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.DataCapture.HttpHeaders.Request = config.Bool(true)
	cfg.Reporting.Endpoint = config.String(srv.URL)

	wrapper := &mockSpanProcessorWrapper{}
	shutdown := InitWithSpanProcessorWrapper(cfg, wrapper, versionInfoAttributes)
	defer shutdown()

	// test wrapper is called for spans created by global trace provider
	_, span, spanEnder := StartSpan(context.Background(), "my_span", nil)
	id1 := span.GetAttributes().GetValue("service.instance.id")
	spanEnder()

	assert.Len(t, id1, 36)

	// my_span and startup spans
	assert.Equal(t, 2, wrapper.onStartCount)
	assert.Equal(t, 2, wrapper.onEndCount)

	// test wrapper is called for spans created by service trace provider
	startSpan, _, err := RegisterServiceWithSpanProcessorWrapper("custom_env", "custom_service", map[string]string{"test1": "val1"}, wrapper,
		versionInfoAttributes)
	if err != nil {
		log.Fatalf("Error while initializing service: %v", err)
	}

	_, serviceSpan, spanEnder := startSpan(context.Background(), "service_span", nil)
	id2 := serviceSpan.GetAttributes().GetValue("service.instance.id")
	assert.Len(t, id2, 36)
	assert.Equal(t, id1, id2)
	spanEnder()

	// service_span, my_span and startup spans
	assert.Equal(t, 3, wrapper.onStartCount)
	assert.Equal(t, 3, wrapper.onEndCount)
}

func TestShouldDisableMetrics(t *testing.T) {
	// Using default values: since zipkin is the default traces exporter turn off metrics
	cfg := config.Load()
	assert.True(t, shouldDisableMetrics(cfg))

	// For OTLP reporting endpoint, turn it on
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP
	assert.False(t, shouldDisableMetrics(cfg))

	cfg = config.Load()
	cfg.Telemetry.MetricsEnabled = config.Bool(false)
	assert.True(t, shouldDisableMetrics(cfg))

	// Set a metrics endpoint
	cfg = config.Load()
	cfg.Reporting.MetricEndpoint = config.String("localhost:4317")
	assert.False(t, shouldDisableMetrics(cfg))
}

func TestShouldUseCustomBatchSpanProcessor(t *testing.T) {
	// Using default values. Should be true
	cfg := config.Load()
	assert.True(t, shouldUseCustomBatchSpanProcessor(cfg))

	cfg.Goagent = nil
	assert.False(t, shouldUseCustomBatchSpanProcessor(cfg))

	cfg.Goagent = &v1.GoAgent{UseCustomBsp: config.Bool(false)}
	assert.False(t, shouldUseCustomBatchSpanProcessor(cfg))

	cfg.Goagent = &v1.GoAgent{}
	assert.False(t, shouldUseCustomBatchSpanProcessor(cfg))

	cfg.Goagent = &v1.GoAgent{UseCustomBsp: config.Bool(true)}
	assert.True(t, shouldUseCustomBatchSpanProcessor(cfg))

	cfg.Telemetry.MetricsEnabled = config.Bool(false)
	assert.False(t, shouldUseCustomBatchSpanProcessor(cfg))
}

func TestConfigFactory(t *testing.T) {
	cfg := config.Load()
	factory := makeConfigFactory(cfg)
	assert.Same(t, cfg, factory())
}
