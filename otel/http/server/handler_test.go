package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traceableai/goagent/otel/internal"
	apitrace "go.opentelemetry.io/otel/api/trace"
)

func TestRequestIsSuccessfullyTraced(t *testing.T) {
	_, flusher := internal.InitTracer()

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("test_res"))
		rw.Write([]byte("ponse_body"))
		rw.Header().Add("request_id", "abc123xyz")
		rw.WriteHeader(202)
	})

	ih := NewHandler(h)

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("test_request_body"))
	r.Header.Add("api_key", "xyz123abc")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	assert.Equal(t, "test_response_body", w.Body.String())

	spans := flusher()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, "GET", spans[0].Name)
	assert.Equal(t, apitrace.SpanKindServer, spans[0].SpanKind)

	for _, kv := range spans[0].Attributes {
		switch kv.Key {
		case "http.status_code":
			assert.Equal(t, "202", kv.Value.AsString())
		case "http.method":
			assert.Equal(t, "GET", kv.Value.AsString())
		case "http.url":
			assert.Equal(t, "http://traceable.ai/foo?user_id=1", kv.Value.AsString())
		case "http.request.header.request_id":
			assert.Equal(t, "abc123xyz", kv.Value.AsString())
		case "http.response.header.api_key":
			assert.Equal(t, "xyz123abc", kv.Value.AsString())
		}
	}
}

func TestRequestAndResponseBodyAreRecordedAccordingly(t *testing.T) {
	_, flusher := internal.InitTracer()

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
			responseContentType:            "application/x-www-form-urlencoded",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"content type and body": {
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "application/x-www-form-urlencoded",
			responseContentType:            "Application/JSON; charset=UTF-8",
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

			span := flusher()[0]
			for _, kv := range span.Attributes {
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

func TestRequestExtractsIncomingHeadersSuccessfully(t *testing.T) {
	_, flusher := internal.InitTracer()

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})

	ih := NewHandler(h)

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("test_request_body"))
	r.Header.Add("X-B3-TraceId", "1f46165474d11ee5836777d85df2cdab")
	r.Header.Add("X-B3-SpanId", "1ee58677d8df2cab")
	r.Header.Add("X-B3-Sampled", "1")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	spans := flusher()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, "1f46165474d11ee5836777d85df2cdab", spans[0].SpanContext.TraceID.String())
	assert.Equal(t, "1ee58677d8df2cab", spans[0].ParentSpanID.String())
}
