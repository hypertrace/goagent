package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

var _ http.Handler = &mockHandler{}

type mockHandler struct {
	baseHandler http.Handler
	spans       []*mock.Span
}

func (h *mockHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	span := mock.NewSpan()
	ctx := mock.ContextWithSpan(context.Background(), span)
	h.spans = append(h.spans, span)
	h.baseHandler.ServeHTTP(rw, r.WithContext(ctx))
}

func TestServerRequestIsSuccessfullyTraced(t *testing.T) {
	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("request_id", "abc123xyz")
		rw.WriteHeader(202)
		rw.Write([]byte("test_res"))
		rw.Write([]byte("ponse_body"))
	})

	wh, _ := WrapHandler(h, mock.SpanFromContext).(*handler)
	wh.dataCaptureConfig = &config.DataCapture{
		HttpHeaders: &config.Message{
			Request:  config.Bool(false),
			Response: config.Bool(false),
		},
		HttpBody: &config.Message{
			Request:  config.Bool(false),
			Response: config.Bool(false),
		},
	}

	ih := &mockHandler{baseHandler: wh}

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("test_request_body"))
	r.Header.Add("api_key", "xyz123abc")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	assert.Equal(t, "test_response_body", w.Body.String())

	assert.Equal(t, 1, len(ih.spans))

	span := ih.spans[0]
	assert.Equal(t, "http://traceable.ai/foo?user_id=1", span.ReadAttribute("http.url").(string))
	assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
}

func TestServerRequestHeadersAreSuccessfullyRecorded(t *testing.T) {
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
		h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add("request_id", "abc123xyz")
			rw.WriteHeader(202)
		})

		wh, _ := WrapHandler(h, mock.SpanFromContext).(*handler)
		ih := &mockHandler{baseHandler: wh}
		wh.dataCaptureConfig = &config.DataCapture{
			HttpHeaders: &config.Message{
				Request:  config.Bool(tCase.captureHTTPHeadersRequestConfig),
				Response: config.Bool(tCase.captureHTTPHeadersResponseConfig),
			},
			HttpBody: &config.Message{
				Request:  config.Bool(false),
				Response: config.Bool(false),
			},
		}

		r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("test_request_body"))
		r.Header.Add("api_key", "xyz123abc")
		w := httptest.NewRecorder()

		ih.ServeHTTP(w, r)

		spans := ih.spans
		assert.Equal(t, 1, len(spans))

		span := spans[0]
		if tCase.captureHTTPHeadersRequestConfig {
			assert.Equal(t, "xyz123abc", span.ReadAttribute("http.request.header.Api_key").(string))
		} else {
			assert.Nil(t, span.ReadAttribute("http.request.header.Api_key"))
		}

		if tCase.captureHTTPHeadersResponseConfig {
			assert.Equal(t, "abc123xyz", span.ReadAttribute("http.response.header.Request_id").(string))
		} else {
			assert.Nil(t, span.ReadAttribute("http.response.header.Request_id"))
		}
	}
}

func TestServerRecordsRequestAndResponseBodyAccordingly(t *testing.T) {
	tCases := map[string]struct {
		captureHTTPBodyConfig          bool
		requestBody                    string
		requestContentType             string
		shouldHaveRecordedRequestBody  bool
		responseBody                   string
		responseContentType            string
		shouldHaveRecordedResponseBody bool
	}{
		"no content type headers and empty body": {
			captureHTTPBodyConfig:          true,
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"no content type headers and non empty body": {
			captureHTTPBodyConfig:          true,
			requestBody:                    "{}",
			responseBody:                   "{}",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
		},
		"content type headers but empty body": {
			captureHTTPBodyConfig:          true,
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
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Add("Content-Type", tCase.responseContentType)
				rw.Header().Add("Content-Type", "charset=UTF-8")
				rw.Write([]byte(tCase.responseBody))
			})

			wh, _ := WrapHandler(h, mock.SpanFromContext).(*handler)
			wh.dataCaptureConfig = &config.DataCapture{
				HttpBody: &config.Message{
					Request:  config.Bool(tCase.captureHTTPBodyConfig),
					Response: config.Bool(tCase.captureHTTPBodyConfig),
				},
				HttpHeaders: &config.Message{
					Request:  config.Bool(false),
					Response: config.Bool(false),
				},
			}
			ih := &mockHandler{baseHandler: wh}

			r, _ := http.NewRequest("GET", "http://traceable.ai/foo", strings.NewReader(tCase.requestBody))
			r.Header.Add("content-type", tCase.requestContentType)

			w := httptest.NewRecorder()

			ih.ServeHTTP(w, r)

			span := ih.spans[0]

			if tCase.shouldHaveRecordedRequestBody {
				assert.Equal(t, tCase.requestBody, span.ReadAttribute("http.request.body").(string))
			} else {
				assert.Nil(t, span.ReadAttribute("http.request.body"))
			}

			if tCase.shouldHaveRecordedResponseBody {
				assert.Equal(t, tCase.responseBody, span.ReadAttribute("http.response.body").(string))
			} else {
				assert.Nil(t, span.ReadAttribute("http.response.body"))
			}
		})
	}
}
