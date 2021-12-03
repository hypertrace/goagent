package opentelemetry

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal/tracetesting"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func TestInitAdditional(t *testing.T) {
	zipkinSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zs := []model.SpanModel{}
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()

		err = json.Unmarshal(b, &zs)
		require.NoError(t, err)

		assert.Equal(t, 1, len(zs))
		assert.Equal(t, "test-span", zs[0].Name)
		assert.Equal(t, "another-name", zs[0].LocalEndpoint.ServiceName)

		w.WriteHeader(http.StatusAccepted)
	}))
	defer zipkinSrv.Close()

	cfg := config.Load()
	cfg.ServiceName = config.String("another-name")
	cfg.Reporting.TraceReporterType = config.TraceReporterType_ZIPKIN
	cfg.Reporting.Endpoint = config.String(zipkinSrv.URL)

	hyperSpanProcessor, shutdown := InitAsAdditional(cfg)
	defer shutdown()

	ctx := context.Background()
	resources, _ := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("one-name"),
		),
	)

	rec := tracetesting.NewRecorder()

	recSpanProcessor := sdktrace.NewSimpleSpanProcessor(RemoveGoAgentAttrs(rec))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(hyperSpanProcessor),
		sdktrace.WithSpanProcessor(recSpanProcessor),
		sdktrace.WithResource(resources),
	)

	defer func() { _ = tp.Shutdown(ctx) }()

	ctx, s := tp.Tracer("test").Start(context.Background(), "test-span")
	s.SetStatus(codes.Ok, "OK")
	s.SetAttributes(
		attribute.String("http.status_code", "200"),
		attribute.String("http.request.header.x-forwarded-for", "1.1.1.1"),
	)
	s.End()

	recSpanProcessor.ForceFlush(context.Background())

	recSpans := rec.Flush()
	assert.Len(t, recSpans, 1)

	recSpan := recSpans[0]
	for _, attr := range recSpan.Attributes() {
		if attr.Key == "http.status_code" {
			assert.Equal(t, "200", attr.Value.AsString())
		}

		if attr.Key == "http.request.header.x-forwarded-for" {
			// This attribute should be filtered out by the RemoveGoAgentAttrs
			t.FailNow()
		}
	}

	for _, attr := range recSpan.Resource().Attributes() {
		if attr.Key == semconv.ServiceNameKey {
			assert.Equal(t, "one-name", attr.Value.AsString())
		}
	}
}
