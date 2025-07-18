package hyperhttp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal/tracetesting"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	"go.opentelemetry.io/otel/propagation"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/trace"
)

func TestClientRequestIsSuccessfullyTraced(t *testing.T) {
	sdkconfig.InitConfig(&config.AgentConfig{
		DataCapture: &config.DataCapture{
			HttpHeaders: &config.Message{
				Request:  config.Bool(true),
				Response: config.Bool(true),
			},
			HttpBody: &config.Message{
				Request:  config.Bool(true),
				Response: config.Bool(true),
			},
			BodyMaxSizeBytes: config.Int32(1024),
			AllowedContentTypes: []*wrapperspb.StringValue{wrapperspb.String("json"),
				wrapperspb.String("x-www-form-urlencoded")},
		},
	})
	defer sdkconfig.ResetConfig()

	_, flusher := tracetesting.InitTracer()

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("content-type", "application/json")
		rw.Header().Set("request_id", "xyz123abc")
		rw.WriteHeader(202)
		rw.Write([]byte(`{"id":123}`))
	}))
	defer srv.Close()

	client := &http.Client{
		Transport: otelhttp.NewTransport(
			WrapTransport(http.DefaultTransport),
		),
	}

	req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(`{"name":"Jacinto"}`))
	req.Header.Set("api_key", "abc123xyz")
	req.Header.Set("content-type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.Equal(t, 202, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"id":123}`, string(resBody))

	spans := flusher()
	assert.Equal(t, 1, len(spans), "unexpected number of spans")

	span := spans[0]
	assert.Equal(t, span.Name(), "HTTP POST")
	assert.Equal(t, trace.SpanKindClient, span.SpanKind())

	attrs := tracetesting.LookupAttributes(span.Attributes())
	// "http.request.method" replaces "http.method"
	assert.Equal(t, "POST", attrs.Get("http.request.method").AsString())
	assert.Equal(t, "abc123xyz", attrs.Get("http.request.header.api_key").AsString())
	assert.Equal(t, `{"name":"Jacinto"}`, attrs.Get("http.request.body").AsString())
	assert.Equal(t, "xyz123abc", attrs.Get("http.response.header.request_id").AsString())
	assert.Equal(t, `{"id":123}`, attrs.Get("http.response.body").AsString())
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
	tracetesting.InitTracer()

	expectedErr := errors.New("roundtrip error")
	client := &http.Client{
		Transport: otelhttp.NewTransport(
			WrapTransport(failingTransport{expectedErr}),
		),
	}

	req, _ := http.NewRequest("POST", "http://test.com", nil)
	_, err := client.Do(req)
	if err == nil {
		t.Errorf("expected error: %v", expectedErr)
	}
}

func TestClientRecordsRequestAndResponseBodyAccordingly(t *testing.T) {
	_, flusher := tracetesting.InitTracer()

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
				Transport: otelhttp.NewTransport(
					WrapTransport(http.DefaultTransport),
				),
			}

			req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(tCase.requestBody))
			req.Header.Set("request_id", "abc123xyz")
			req.Header.Set("content-type", tCase.requestContentType)
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			assert.Equal(t, 202, res.StatusCode)

			_, err = io.ReadAll(res.Body)
			assert.Nil(t, err)

			span := flusher()[0]

			var attrs tracetesting.LookupAttributes = tracetesting.LookupAttributes(span.Attributes())
			if tCase.shouldHaveRecordedRequestBody {
				assert.Equal(t, tCase.requestBody, attrs.Get("http.request.body").AsString())
			} else {
				assert.Equal(t, "", attrs.Get("http.request.body").AsString())
			}

			if tCase.shouldHaveRecordedResponseBody {
				assert.Equal(t, tCase.responseBody, attrs.Get("http.response.body").AsString())
			} else {
				assert.Equal(t, "", attrs.Get("http.response.body").AsString())
			}
		})
	}
}

func TestTransportRequestInjectsHeadersSuccessfully(t *testing.T) {
	tracer, _ := tracetesting.InitTracer()

	ctx, span := tracer.Start(context.Background(), "test")
	defer span.End()

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// We make sure the context is being injected.
		ctx := b3.New().Extract(context.Background(), propagation.HeaderCarrier(req.Header))
		_, extractedSpan := tracer.Start(ctx, "test2")
		defer extractedSpan.End()
		assert.Equal(t, span.SpanContext().TraceID().String(), extractedSpan.SpanContext().TraceID().String())
	}))
	defer srv.Close()

	client := &http.Client{
		Transport: otelhttp.NewTransport(
			WrapTransport(http.DefaultTransport),
		),
	}

	req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(`{"name":"Jacinto"}`))
	req.Header.Set("request_id", "abc123xyz")
	req.Header.Set("content-type", "application/json")
	_, err := client.Do(req.WithContext(ctx))
	assert.Nil(t, err)
}
