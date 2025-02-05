package opentelemetry

import (
	"context"
	"fmt"
	v1 "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
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

	_, _, err := RegisterService("custom_service", map[string]string{"test1": "val1"})
	if err != nil {
		log.Fatalf("Error while initializing service: %v", err)
	}
}

func TestInitDisabledAgent(t *testing.T) {
	cfg := config.Load()
	cfg.Enabled = config.Bool(false)
	shutdown := Init(cfg)
	defer shutdown()

	startSpan, tp, err := RegisterService("test_service", nil)
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

	startSpan, tp, err := RegisterService("custom_service", map[string]string{"test1": "val1"})
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
	_, tp, err := RegisterService("custom_service", map[string]string{"test1": "val1"})
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

	startServiceSpan, tp, err := RegisterService("custom_service", map[string]string{"test1": "val1"})
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

	startServiceSpan, _, err := RegisterService("custom_service", map[string]string{"test1": "val1"})
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
	startSpan, _, err := RegisterServiceWithSpanProcessorWrapper("custom_service", map[string]string{"test1": "val1"}, wrapper,
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

type MockTraceService struct {
	coltracepb.UnimplementedTraceServiceServer
	metadataCh chan metadata.MD
}

func (m *MockTraceService) Export(ctx context.Context, req *coltracepb.ExportTraceServiceRequest) (*coltracepb.ExportTraceServiceResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println("Received Metadata:", md)

	m.metadataCh <- md
	return &coltracepb.ExportTraceServiceResponse{}, nil
}

func waitForServer(address string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("server did not start within %v", timeout)
}

func TestMakeExporterFactory_Headers_WithMockGRPCServer(t *testing.T) {
	metadataCh := make(chan metadata.MD, 1)

	listener, err := net.Listen("tcp", "127.0.0.1:49999")
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	mockTraceService := &MockTraceService{metadataCh: metadataCh}
	coltracepb.RegisterTraceServiceServer(grpcServer, mockTraceService)

	go func() {
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	defer grpcServer.Stop()

	require.NoError(t, waitForServer("127.0.0.1:49999", 5*time.Second))

	cfg := &v1.AgentConfig{
		Reporting: &v1.Reporting{
			Token:             config.String("test-token"),
			TraceReporterType: v1.TraceReporterType_OTLP,
			Endpoint:          config.String("127.0.0.1:49999"),
		},
	}

	exporterFactory := makeExporterFactory(cfg)

	// Create the exporter
	exporter, err := exporterFactory()
	require.NoError(t, err)
	require.NotNil(t, exporter)

	tp := sdktrace.NewTracerProvider()
	_, span := tp.Tracer("test-tracer").Start(context.Background(), "test-span")
	span.End()

	err = exporter.ExportSpans(context.Background(), []sdktrace.ReadOnlySpan{span.(sdktrace.ReadOnlySpan)})
	assert.NoError(t, err)

	select {
	case capturedMetadata := <-metadataCh:
		require.NotNil(t, capturedMetadata)
		assert.Equal(t, []string{"test-token"}, capturedMetadata["todo-traceable-agent-token"])
	case <-time.After(1 * time.Second):
		t.Fatal("Metadata was not captured")
	}

	if exporter != nil {
		_ = exporter.Shutdown(context.Background())
	}
}

func TestMakeExporterFactory_Headers_ZipkinAndOTLPHTTP(t *testing.T) {
	// Channel to capture headers
	headersCh := make(chan http.Header, 1)

	// Mock HTTP server to capture requests
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headersCh <- r.Header              // Capture headers
		w.WriteHeader(http.StatusAccepted) // zipkin spec returns 202, otlp http doesnt care but zipkin exporter will fail otherwise
	}))
	defer mockServer.Close()

	testCases := []struct {
		reporterType v1.TraceReporterType
	}{
		{v1.TraceReporterType_ZIPKIN},
		{v1.TraceReporterType_OTLP_HTTP},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing %v", tc.reporterType), func(t *testing.T) {
			endpoint := mockServer.URL
			// otlp http doesn't want scheme, zipkin does
			if tc.reporterType == v1.TraceReporterType_OTLP_HTTP {
				endpoint = removeProtocolPrefixForOTLP(endpoint)
			}

			cfg := &v1.AgentConfig{
				Reporting: &v1.Reporting{
					Token:             config.String("test-token"),
					TraceReporterType: tc.reporterType,
					Endpoint:          config.String(endpoint),
				},
			}

			exporterFactory := makeExporterFactory(cfg)

			exporter, err := exporterFactory()
			require.NoError(t, err)
			require.NotNil(t, exporter)

			tp := sdktrace.NewTracerProvider()
			_, span := tp.Tracer("test-tracer").Start(context.Background(), "test-span")
			span.End()

			err = exporter.ExportSpans(context.Background(), []sdktrace.ReadOnlySpan{span.(sdktrace.ReadOnlySpan)})
			assert.NoError(t, err)

			select {
			case capturedHeaders := <-headersCh:
				require.NotNil(t, capturedHeaders)
				assert.Equal(t, []string{"test-token"}, capturedHeaders.Values("TODO-traceable-agent-token"))
			case <-time.After(2 * time.Second):
				t.Fatal("Metadata was not captured")
			}

			if exporter != nil {
				_ = exporter.Shutdown(context.Background())
			}
		})
	}
}
