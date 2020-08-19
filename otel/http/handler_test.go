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

type recorder struct {
	spans []*trace.SpanData
}

// ExportSpan records a span
func (r *recorder) ExportSpan(ctx context.Context, s *trace.SpanData) {
	r.spans = append(r.spans, s)
}

// Flush returns the current recorded spans and reset the recordings
func (r *recorder) Flush() []*trace.SpanData {
	spans := r.spans
	r.spans = nil
	return spans
}

func initTracer() func() []*trace.SpanData {
	exporter := &recorder{}
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
		rw.Write([]byte("test_res"))
		rw.Write([]byte("ponse_body"))
		rw.Header().Add("request_id", "abc123xyz")
		rw.WriteHeader(202)
	})

	ih := NewHandler(h)

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo", strings.NewReader("test_request_body"))
	r.Header.Add("api_key", "xyz123abc")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	assert.Equal(t, "test_response_body", w.Body.String())

	spans := flusher()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, "GET", spans[0].Name)

	for _, kv := range spans[0].Attributes {
		switch kv.Key {
		case "http.status_code":
			assert.Equal(t, "202", kv.Value.AsString())
		case "http.method":
			assert.Equal(t, "GET", kv.Value.AsString())
		case "http.path":
			assert.Equal(t, "/foo", kv.Value.AsString())
		case "http.request.header.request_id":
			assert.Equal(t, "abc123xyz", kv.Value.AsString())
		case "http.response.header.api_key":
			assert.Equal(t, "xyz123abc", kv.Value.AsString())
		}
	}
}

func TestRequestAndResponseBodyAreRecordedAccordingly(t *testing.T) {
	flusher := initTracer()

	tCases := map[string]struct {
		requestBody                    string
		requestContentType             string
		shouldHaveRecordedRequestBody  bool
		responseBody                   string
		responseContentType            string
		shouldHaveRecordedResponseBody bool
	}{
		"no content type headers and empty body": {
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"no content type headers and non empty body": {
			requestBody:                    "{}",
			responseBody:                   "{}",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"content type headers but empty body": {
			requestContentType:             "application/json",
			responseContentType:            "application/json",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"content type and body": {
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "application/json",
			responseContentType:            "application/json",
			shouldHaveRecordedRequestBody:  true,
			shouldHaveRecordedResponseBody: true,
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Write([]byte(tCase.responseBody))
				rw.Header().Add("Content-Type", tCase.responseContentType)
			})

			ih := NewHandler(h)

			r, _ := http.NewRequest("GET", "http://traceable.ai/foo", strings.NewReader(tCase.requestBody))
			r.Header.Add("content-type", tCase.requestContentType)

			w := httptest.NewRecorder()

			ih.ServeHTTP(w, r)

			spans := flusher()
			assert.Equal(t, 1, len(spans))
			assert.Equal(t, "GET", spans[0].Name)

			for _, kv := range spans[0].Attributes {
				switch kv.Key {
				case "http.request.body":
					if tCase.shouldHaveRecordedRequestBody {
						assert.Equal(t, tCase.requestBody, kv.Value.AsString())
					} else {
						t.Errorf("unexpected request body recording")
					}
				case "http.response.body":
					if tCase.shouldHaveRecordedResponseBody {
						assert.Equal(t, tCase.responseBody, kv.Value.AsString())
					} else {
						t.Errorf("unexpected response body recording")
					}
				}
			}
		})
	}
}
