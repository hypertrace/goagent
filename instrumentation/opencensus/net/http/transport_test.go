package http

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hypertrace/goagent/instrumentation/opencensus/internal"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
)

func TestClientRequestIsSuccessfullyTraced(t *testing.T) {
	flusher := internal.InitTracer()

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("content-type", "application/json")
		rw.Header().Set("request_id", "xyz123abc")
		rw.WriteHeader(202)
		rw.Write([]byte(`{"id":123}`))
	}))
	defer srv.Close()

	client := &http.Client{
		Transport: &ochttp.Transport{
			Base: WrapTransport(http.DefaultTransport),
		},
	}

	req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(`{"name":"Jacinto"}`))
	req.Header.Set("api_key", "abc123xyz")
	req.Header.Set("content-type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.Equal(t, 202, res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"id":123}`, string(resBody))

	spans := flusher()
	assert.Equal(t, 1, len(spans), "unexpected number of spans")

	span := spans[0]
	assert.Equal(t, trace.SpanKindClient, span.SpanKind)

	assert.Equal(t, "POST", span.Attributes["http.method"].(string))
	assert.Equal(t, "abc123xyz", span.Attributes["http.request.header.Api_key"].(string))
	assert.Equal(t, `{"name":"Jacinto"}`, span.Attributes["http.request.body"].(string))
	assert.Equal(t, "xyz123abc", span.Attributes["http.response.header.Request_id"].(string))
	assert.Equal(t, `{"id":123}`, span.Attributes["http.response.body"].(string))
}

type failingTransport struct {
	err error
}

func (t failingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	if t.err == nil {
		log.Fatal("missing error in failing transport")
	}
	return nil, t.err
}

func TestClientFailureRequestIsSuccessfullyTraced(t *testing.T) {
	internal.InitTracer()

	expectedErr := errors.New("roundtrip error")
	client := &http.Client{
		Transport: &ochttp.Transport{
			Base: WrapTransport(failingTransport{expectedErr}),
		},
	}

	req, _ := http.NewRequest("POST", "http://test.com", nil)
	_, err := client.Do(req)
	if err == nil {
		t.Errorf("expected error: %v", expectedErr)
	}
}

func TestClientRecordsRequestAndResponseBodyAccordingly(t *testing.T) {
	flusher := internal.InitTracer()

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
			responseContentType:            "Application/JSON",
			shouldHaveRecordedRequestBody:  true,
			shouldHaveRecordedResponseBody: true,
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Header().Add("Content-Type", tCase.responseContentType)
				rw.Header().Add("Content-Type", "charset=UTF-8")
				rw.WriteHeader(202)
				rw.Write([]byte(tCase.responseBody))
			}))
			defer srv.Close()

			client := &http.Client{
				Transport: &ochttp.Transport{
					Base: WrapTransport(http.DefaultTransport),
				},
			}

			req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(tCase.requestBody))
			req.Header.Set("request_id", "abc123xyz")
			req.Header.Set("content-type", tCase.requestContentType)
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			assert.Equal(t, 202, res.StatusCode)

			_, err = ioutil.ReadAll(res.Body)
			assert.Nil(t, err)

			span := flusher()[0]

			if tCase.shouldHaveRecordedRequestBody {
				assert.Equal(t, tCase.requestBody, span.Attributes["http.request.body"].(string))
			} else {
				_, found := span.Attributes["http.request.body"]
				assert.False(t, found)
			}

			if tCase.shouldHaveRecordedResponseBody {
				assert.Equal(t, tCase.responseBody, span.Attributes["http.response.body"].(string))
			} else {
				_, found := span.Attributes["http.response.body"]
				assert.False(t, found)
			}
		})
	}
}

func TestTransportRequestInjectsHeadersSuccessfully(t *testing.T) {
	internal.InitTracer()
	ctx, span := trace.StartSpan(context.Background(), "test")
	defer span.End()

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// We make sure the context is being injected.
		propagator := &b3.HTTPFormat{}
		sc, ok := propagator.SpanContextFromRequest(req)
		assert.True(t, ok)
		_, extractedSpan := trace.StartSpanWithRemoteParent(ctx, "test2", sc)
		assert.Equal(t, span.SpanContext().TraceID, extractedSpan.SpanContext().TraceID)
	}))
	defer srv.Close()

	client := &http.Client{
		Transport: &ochttp.Transport{
			Base: WrapTransport(http.DefaultTransport),
		},
	}

	req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(`{"name":"Jacinto"}`))
	req.Header.Set("request_id", "abc123xyz")
	req.Header.Set("content-type", "application/json")
	_, err := client.Do(req.WithContext(ctx))
	assert.Nil(t, err)
}
