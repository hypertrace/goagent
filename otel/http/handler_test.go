package http

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/standard"
	"go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Recorder struct {
	spans []*trace.SpanData
}

func (r *Recorder) ExportSpan(ctx context.Context, s *trace.SpanData) {
	r.spans = append(r.spans, s)
}

func (r *Recorder) Flush() []*trace.SpanData {
	return r.spans
}

func initTracer() func() []*trace.SpanData {
	exporter := &Recorder{}
	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.New(standard.ServiceNameKey.String("ExampleService"))))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
	return func() []*trace.SpanData {
		return exporter.Flush()
	}
}

func TestRequestIsSuccessfullyTraced(t *testing.T) {
	flusher := initTracer()

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("test_response_body"))
	})

	ih := NewHandler(h)

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo", strings.NewReader("test_request_body"))
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	assert.Equal(t, "test_response_body", w.Body.String())

	spans := flusher()

	assert.Equal(t, 1, len(spans))
	assert.Equal(t, "GET", spans[0].Name)

	for _, kv := range spans[0].Attributes {
		switch kv.Key {
		case "http.method":
			assert.Equal(t, "GET", kv.Value.AsString())
		case "http.path":
			assert.Equal(t, "/foo", kv.Value.AsString())
		case "http.request.body":
			assert.Equal(t, "test_request_body", kv.Value.AsString())
		case "http.response.body":
			assert.Equal(t, "test_response_body", kv.Value.AsString())
		}
	}
}
