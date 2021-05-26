package opentelemetry

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.opentelemetry.io/otel/propagation"

	"github.com/hypertrace/goagent/config"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
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
	cfg.Reporting.Endpoint = config.String("http://api.traceable.ai:9411/api/v2/spans")
	cfg.Reporting.TraceReporterType = config.TraceReporterType_ZIPKIN

	shutdown := Init(cfg)

	_, _, err := RegisterService("custom_service", map[string]string{"test1": "val1"})
	if err != nil {
		log.Fatalf("Error while initializing service: %v", err)
	}

	defer shutdown()
}

func TestShutdownFlushesAllSpans(t *testing.T) {
	requestIsReceived := false
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		requestIsReceived = true
	}))
	defer srv.Close()

	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.Reporting.Endpoint = config.String(srv.URL)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_ZIPKIN

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
	}))
	defer srv.Close()

	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.Reporting.Endpoint = config.String(srv.URL)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_ZIPKIN

	// By doing this we make sure a batching isn't happening
	batchTimeout = time.Duration(10) * time.Second

	shutdown := Init(cfg)
	assert.True(t, initialized)
	assert.Equal(t, 0, len(traceProviders))

	_, _, spanEnder := StartSpan(context.Background(), "example_span", nil)
	spanEnder()

	startServiceSpan, _, err := RegisterService("custom_service", map[string]string{"test1": "val1"})
	assert.NoError(t, err)
	assert.NotNil(t, startServiceSpan)
	assert.True(t, initialized)
	assert.Equal(t, 1, len(traceProviders))

	_, _, serviceSpanEnder := startServiceSpan(context.Background(), "my_span", nil)
	serviceSpanEnder()

	assert.Equal(t, 0, count)
	shutdown()
	assert.Equal(t, 2, count)
	assert.Equal(t, 0, len(traceProviders))
	fmt.Println("Count: ", count)
}

func TestMultipleTraceProvidersCallAfterShutdown(t *testing.T) {
	requestIsReceived := false
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		requestIsReceived = true
	}))
	defer srv.Close()

	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.Reporting.Endpoint = config.String(srv.URL)
	cfg.Reporting.TraceReporterType = config.TraceReporterType_ZIPKIN

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
	cfg.PropagationFormats = []config.PropagationFormat{
		config.PropagationFormat_B3,
		config.PropagationFormat_TRACECONTEXT,
	}
	Init(cfg)
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
	Init(cfg)
}
