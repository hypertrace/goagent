package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hypertrace/goagent/config"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
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

func TestMain(m *testing.M) {
	sdkconfig.InitConfig(config.AgentConfig{
		DataCapture: &config.DataCapture{
			HTTPHeaders: &config.Message{
				Request:  config.BoolVal(true),
				Response: config.BoolVal(true),
			},
			HTTPBody: &config.Message{
				Request:  config.BoolVal(true),
				Response: config.BoolVal(true),
			},
		},
	})
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
		HTTPHeaders: &config.Message{
			Request:  config.BoolVal(true),
			Response: config.BoolVal(true),
		},
	}
	ih := &mockHandler{baseHandler: wh}

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("test_request_body"))
	r.Header.Add("api_key", "xyz123abc")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	assert.Equal(t, "test_response_body", w.Body.String())

	spans := ih.spans
	assert.Equal(t, 1, len(spans))

	assert.Equal(t, "http://traceable.ai/foo?user_id=1", spans[0].Attributes["http.url"].(string))
	assert.Equal(t, "xyz123abc", spans[0].Attributes["http.request.header.Api_key"].(string))
	assert.Equal(t, "abc123xyz", spans[0].Attributes["http.response.header.Request_id"].(string))
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
				HTTPBody: &config.Message{
					Request:  &tCase.captureHTTPBodyConfig,
					Response: &tCase.captureHTTPBodyConfig,
				},
			}
			ih := &mockHandler{baseHandler: wh}

			r, _ := http.NewRequest("GET", "http://traceable.ai/foo", strings.NewReader(tCase.requestBody))
			r.Header.Add("content-type", tCase.requestContentType)

			w := httptest.NewRecorder()

			ih.ServeHTTP(w, r)

			span := ih.spans[0]

			if tCase.shouldHaveRecordedRequestBody {
				assert.Equal(t, tCase.requestBody, span.Attributes["http.request.body"].(string))
			}

			if tCase.shouldHaveRecordedResponseBody {
				assert.Equal(t, tCase.responseBody, span.Attributes["http.response.body"].(string))
			}
		})
	}
}
