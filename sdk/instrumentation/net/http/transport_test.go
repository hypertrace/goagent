package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter/result"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

type mockTransport struct {
	baseRoundTripper http.RoundTripper
	spans            []*mock.Span
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	span := mock.NewSpan()
	ctx := mock.ContextWithSpan(context.Background(), span)
	t.spans = append(t.spans, span)
	return t.baseRoundTripper.RoundTrip(req.WithContext(ctx))
}

func TestClientRequestIsSuccessfullyTraced(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(202)
		rw.Write([]byte(`{"id":123}`))
	}))
	defer srv.Close()

	rt := WrapTransport(http.DefaultTransport, mock.SpanFromContext, &Options{}, map[string]string{"foo": "bar"}).(*roundTripper)
	rt.dataCaptureConfig = &config.DataCapture{
		HttpHeaders: &config.Message{
			Request:  config.Bool(false),
			Response: config.Bool(false),
		},
		HttpBody: &config.Message{
			Request:  config.Bool(false),
			Response: config.Bool(false),
		},
		BodyMaxSizeBytes: config.Int32(1000),
	}

	tr := &mockTransport{
		baseRoundTripper: rt,
	}
	client := &http.Client{
		Transport: tr,
	}

	req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(`{"name":"Jacinto"}`))
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.Equal(t, 202, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"id":123}`, string(resBody))

	spans := tr.spans
	assert.Equal(t, 1, len(spans), "unexpected number of spans")

	span := spans[0]

	_ = span.ReadAttribute("container_id") // needed in containarized envs
	// custom attribute
	assert.Equal(t, "bar", span.ReadAttribute("foo").(string))
	// We make sure we read all attributes and covered them with tests
	assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
}

func TestClientRequestHeadersAreCapturedAccordingly(t *testing.T) {
	tCases := []struct {
		captureHTTPHeadersRequestConfig  bool
		captureHTTPHeadersResponseConfig bool
	}{
		{true, true},
		{true, false},
		{false, true},
		{false, false},
	}
	for _, tCase := range tCases {

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("content-type", "application/json")
			rw.Header().Set("request_id", "xyz123abc")
			rw.WriteHeader(202)
			rw.Write([]byte(`{"id":123}`))
		}))
		defer srv.Close()

		rt := WrapTransport(http.DefaultTransport, mock.SpanFromContext, &Options{}, map[string]string{"foo": "bar"}).(*roundTripper)
		rt.dataCaptureConfig = &config.DataCapture{
			HttpHeaders: &config.Message{
				Request:  config.Bool(tCase.captureHTTPHeadersRequestConfig),
				Response: config.Bool(tCase.captureHTTPHeadersResponseConfig),
			},
			HttpBody: &config.Message{
				Request:  config.Bool(false),
				Response: config.Bool(false),
			},
			BodyMaxSizeBytes: config.Int32(1000),
		}

		tr := &mockTransport{
			baseRoundTripper: rt,
		}
		client := &http.Client{
			Transport: tr,
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

		spans := tr.spans
		assert.Equal(t, 1, len(spans), "unexpected number of spans")

		span := spans[0]
		if tCase.captureHTTPHeadersRequestConfig {
			assert.Equal(t, "abc123xyz", span.ReadAttribute("http.request.header.api_key").(string))
		} else {
			assert.Nil(t, span.ReadAttribute("http.request.header.Api_key"))
		}

		if tCase.captureHTTPHeadersResponseConfig {
			assert.Equal(t, "xyz123abc", span.ReadAttribute("http.response.header.request_id").(string))
		} else {
			assert.Nil(t, span.ReadAttribute("http.response.header.request_id"))
		}
		assert.Equal(t, "bar", span.ReadAttribute("foo").(string))
	}
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
	expectedErr := errors.New("roundtrip error")
	client := &http.Client{
		Transport: &mockTransport{
			baseRoundTripper: WrapTransport(failingTransport{expectedErr}, mock.SpanFromContext, &Options{}, map[string]string{}),
		},
	}

	req, _ := http.NewRequest("POST", "http://test.com", nil)
	_, err := client.Do(req)
	if err == nil {
		t.Errorf("expected error: %v", expectedErr)
	}
}

func TestClientRecordsRequestAndResponseBodyAccordingly(t *testing.T) {
	tCases := map[string]struct {
		captureHTTPBodyConfig          bool
		requestBody                    interface{}
		requestContentType             string
		shouldHaveRecordedRequestBody  bool
		shouldBase64EncodeRequestBody  bool
		responseBody                   string
		responseContentType            string
		shouldHaveRecordedResponseBody bool
		shouldBase64EncodeResponseBody bool
	}{
		"no content type headers and empty body": {
			captureHTTPBodyConfig: true,

			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"nil body": {
			captureHTTPBodyConfig:          true,
			requestBody:                    nil,
			requestContentType:             "application/json",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"no content type headers and non empty body": {
			captureHTTPBodyConfig: true,

			requestBody:                    "{}",
			responseBody:                   "{}",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"content type headers but empty body": {
			captureHTTPBodyConfig: true,

			requestContentType:             "application/json",
			responseContentType:            "application/x-www-form-urlencoded",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"content type and body with config enabled": {
			captureHTTPBodyConfig:          true,
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "application/x-www-form-urlencoded",
			responseContentType:            "Application/JSON",
			shouldHaveRecordedRequestBody:  true,
			shouldHaveRecordedResponseBody: true,
		},
		"content type and body but config disabled": {
			captureHTTPBodyConfig:          false,
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "application/x-www-form-urlencoded",
			responseContentType:            "Application/JSON",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"request multipart/form-data content type and body config enabled": {
			captureHTTPBodyConfig:          true,
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "multipart/form-data",
			responseContentType:            "application/json",
			shouldHaveRecordedRequestBody:  true,
			shouldHaveRecordedResponseBody: true,
			shouldBase64EncodeRequestBody:  true,
			shouldBase64EncodeResponseBody: false,
		},
		"response multipart/form-data content type and body config enabled": {
			captureHTTPBodyConfig:          true,
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "application/json",
			responseContentType:            "multipart/form-data",
			shouldHaveRecordedRequestBody:  true,
			shouldHaveRecordedResponseBody: true,
			shouldBase64EncodeRequestBody:  false,
			shouldBase64EncodeResponseBody: true,
		},
		"request and response multipart/form-data content type and body config enabled": {
			captureHTTPBodyConfig:          true,
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "multipart/form-data",
			responseContentType:            "multipart/form-data",
			shouldHaveRecordedRequestBody:  true,
			shouldHaveRecordedResponseBody: true,
			shouldBase64EncodeRequestBody:  true,
			shouldBase64EncodeResponseBody: true,
		},
		"request and response multipart/form-data content type and body config disabled": {
			captureHTTPBodyConfig: false,
			requestBody:           "test_request_body",
			responseBody:          "test_response_body",
			requestContentType:    "multipart/form-data",
			responseContentType:   "multipart/form-data",
		},
		"request and response multipart/form-data content type not allowed and body config enabled ": {
			captureHTTPBodyConfig: true,
			requestBody:           "test_request_body",
			responseBody:          "test_response_body",
			requestContentType:    "multipart/form-data",
			responseContentType:   "multipart/form-data",
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

			rt := WrapTransport(http.DefaultTransport, mock.SpanFromContext, &Options{}, map[string]string{}).(*roundTripper)
			rt.dataCaptureConfig = &config.DataCapture{
				HttpBody: &config.Message{
					Request:  config.Bool(tCase.captureHTTPBodyConfig),
					Response: config.Bool(tCase.captureHTTPBodyConfig),
				},
				HttpHeaders: &config.Message{
					Request:  config.Bool(false),
					Response: config.Bool(false),
				},
				BodyMaxSizeBytes: config.Int32(1000),
			}
			defaultAllowedContentTypes := internalconfig.GetConfig().DataCapture.AllowedContentTypes
			// add multipart/form-data to allowed content-types
			if tCase.shouldBase64EncodeRequestBody || tCase.shouldBase64EncodeResponseBody {
				internalconfig.GetConfig().DataCapture.AllowedContentTypes = append(internalconfig.GetConfig().DataCapture.AllowedContentTypes,
					config.String("multipart/form-data"))
			}

			tr := &mockTransport{
				baseRoundTripper: rt,
			}

			client := &http.Client{
				Transport: tr,
			}

			var reqBody io.Reader
			if tCase.requestBody != nil {
				reqBody = bytes.NewBufferString(tCase.requestBody.(string))
			} else {
				reqBody = nil
			}
			req, _ := http.NewRequest("POST", srv.URL, reqBody)
			req.Header.Set("request_id", "abc123xyz")
			req.Header.Set("content-type", tCase.requestContentType)
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			assert.Equal(t, 202, res.StatusCode)

			_, err = io.ReadAll(res.Body)
			assert.Nil(t, err)

			span := tr.spans[0]
			if tCase.shouldHaveRecordedRequestBody {
				if tCase.shouldBase64EncodeRequestBody {
					assert.Nil(t, span.ReadAttribute("http.request.body"))
					assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte(tCase.requestBody.(string))),
						span.ReadAttribute("http.request.body.base64").(string))
				} else {
					assert.Equal(t, tCase.requestBody.(string), span.ReadAttribute("http.request.body").(string))
					assert.Nil(t, span.ReadAttribute("http.request.body.base64"))
				}
			} else {
				assert.Nil(t, span.ReadAttribute("http.request.body"))
				assert.Nil(t, span.ReadAttribute("http.request.body.base64"))
			}

			if tCase.shouldHaveRecordedResponseBody {
				if tCase.shouldBase64EncodeResponseBody {
					assert.Nil(t, span.ReadAttribute("http.response.body"))
					assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte(tCase.responseBody)),
						span.ReadAttribute("http.response.body.base64").(string))
				} else {
					assert.Equal(t, tCase.responseBody, span.ReadAttribute("http.response.body").(string))
					assert.Nil(t, span.ReadAttribute("http.response.body.base64"))
				}
			} else {
				assert.Nil(t, span.ReadAttribute("http.response.body"))
				assert.Nil(t, span.ReadAttribute("http.response.body.base64"))
			}

			// reset allowed content types config
			internalconfig.GetConfig().DataCapture.AllowedContentTypes = defaultAllowedContentTypes
		})
	}
}

func TestFilter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(202)
		rw.Write([]byte(`{"id":123}`))
	}))
	defer srv.Close()

	dcCfg := &config.DataCapture{
		HttpHeaders: &config.Message{
			Request:  config.Bool(false),
			Response: config.Bool(false),
		},
		HttpBody: &config.Message{
			Request:  config.Bool(false),
			Response: config.Bool(false),
		},
		BodyMaxSizeBytes: config.Int32(1000),
	}

	tests := []struct {
		name  string
		block bool
	}{
		{
			name:  "blocking enabled",
			block: true,
		},
		{
			name:  "blocking disabled",
			block: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := mock.Filter{
				Evaluator: func(span sdk.Span) result.FilterResult {
					span.SetAttribute("filter.evaluated", true)
					return result.FilterResult{
						Block:              tt.block,
						ResponseStatusCode: 403,
						ResponseMessage:    "Access Denied",
					}
				},
			}
			rt := WrapTransport(http.DefaultTransport, mock.SpanFromContext, &Options{
				Filter: filter,
			}, map[string]string{"foo": "bar"}).(*roundTripper)
			rt.dataCaptureConfig = dcCfg

			tr := &mockTransport{
				baseRoundTripper: rt,
			}
			client := &http.Client{
				Transport: tr,
			}

			req, _ := http.NewRequest("POST", srv.URL, bytes.NewBufferString(`{"name":"Jacinto"}`))
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.block {
				assert.Equal(t, 202, res.StatusCode)
				resBody, err := io.ReadAll(res.Body)
				assert.Nil(t, err)
				assert.Equal(t, `{"id":123}`, string(resBody))

				spans := tr.spans
				assert.Equal(t, 1, len(spans), "unexpected number of spans")

				span := spans[0]

				_ = span.ReadAttribute("container_id") // needed in containarized envs
				// custom attribute
				assert.Equal(t, "bar", span.ReadAttribute("foo").(string))
				assert.True(t, span.ReadAttribute("filter.evaluated").(bool))
				// We make sure we read all attributes and covered them with tests
				assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
			} else {
				assert.Equal(t, 403, res.StatusCode)
				resBody, err := io.ReadAll(res.Body)
				assert.Nil(t, err)
				assert.Equal(t, `Access Denied`, string(resBody))

				spans := tr.spans
				assert.Equal(t, 1, len(spans), "unexpected number of spans")

				span := spans[0]

				_ = span.ReadAttribute("container_id") // needed in containarized envs
				// custom attribute
				assert.Equal(t, "bar", span.ReadAttribute("foo").(string))
				assert.Equal(t, int32(403), span.ReadAttribute("http.status_code").(int32))
				assert.True(t, span.ReadAttribute("filter.evaluated").(bool))
				// We make sure we read all attributes and covered them with tests
				assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)

			}
		})
	}

}
