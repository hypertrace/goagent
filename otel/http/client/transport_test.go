package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	internal "github.com/traceableai/goagent/otel/internal"
	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	apitrace "go.opentelemetry.io/otel/api/trace"
)

func TestTransportRecordsRequestAndResponseBody(t *testing.T) {
	_, flusher := internal.InitTracer()

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(202)
		rw.Write([]byte(`{"id":123}`))
		rw.Header().Set("content-type", "application/json")
		rw.Header().Set("response_id", "xyz123abc")
	}))
	defer srv.Close()

	client := &http.Client{
		Transport: otelhttp.NewTransport(
			Wrap(http.DefaultTransport),
		),
	}

	req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(`{"name":"Jacinto"}`))
	req.Header.Set("request_id", "abc123xyz")
	req.Header.Set("content-type", "application/json")
	res, err := client.Do(req)

	assert.Equal(t, 202, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"id":123}`, string(resBody))

	spans := flusher()
	assert.Equal(t, 1, len(spans), "unexpected number of spans")

	span := spans[0]
	assert.Equal(t, span.Name, "POST")
	assert.Equal(t, span.SpanKind, apitrace.SpanKindClient)
	for _, kv := range span.Attributes {
		switch kv.Key {
		case "http.method":
			assert.Equal(t, "POST", kv.Value.AsString())
		case "http.request.header.request_id":
			assert.Equal(t, "abc123xyz", kv.Value.AsString())
		case "http.request.header.response_id":
			assert.Equal(t, "xyz123abc", kv.Value.AsString())
		case "http.request.body":
			assert.Equal(t, `{"name":"Jacinto"}`, kv.Value.AsString())
		case "http.response.body":
			assert.Equal(t, `{"id":123}`, kv.Value.AsString())
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
			srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(202)
				rw.Write([]byte(tCase.responseBody))
				rw.Header().Add("Content-Type", tCase.responseContentType)
			}))
			defer srv.Close()

			client := &http.Client{
				Transport: otelhttp.NewTransport(
					Wrap(http.DefaultTransport),
				),
			}

			req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(tCase.requestBody))
			req.Header.Set("request_id", "abc123xyz")
			req.Header.Set("content-type", tCase.requestContentType)
			res, err := client.Do(req)

			assert.Equal(t, 202, res.StatusCode)

			_, err = ioutil.ReadAll(res.Body)
			assert.Nil(t, err)

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

func TestRequestInjectsHeadersSuccessfully(t *testing.T) {
	tracer, _ := internal.InitTracer()

	ctx, span := tracer.Start(context.Background(), "test")
	defer span.End()

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// We make sure the context is being injected.
		ctx := propagation.ExtractHTTP(context.Background(), global.Propagators(), req.Header)
		_, extractedSpan := tracer.Start(ctx, "test2")
		assert.Equal(t, span.SpanContext().TraceID, extractedSpan.SpanContext().TraceID)
	}))
	defer srv.Close()

	client := &http.Client{
		Transport: otelhttp.NewTransport(
			Wrap(http.DefaultTransport),
		),
	}

	req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(`{"name":"Jacinto"}`))
	req.Header.Set("request_id", "abc123xyz")
	req.Header.Set("content-type", "application/json")
	_, err := client.Do(req.WithContext(ctx))
	assert.Nil(t, err)
}
