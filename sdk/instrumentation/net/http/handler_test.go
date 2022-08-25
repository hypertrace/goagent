package http

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter/result"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

var emptyTestConfig = &config.DataCapture{
	HttpHeaders: &config.Message{
		Request:  config.Bool(false),
		Response: config.Bool(false),
	},
	HttpBody: &config.Message{
		Request:  config.Bool(false),
		Response: config.Bool(false),
	},
	BodyMaxSizeBytes:           config.Int32(1000),
	BodyMaxProcessingSizeBytes: config.Int32(1000),
}

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

func TestServerRequestWithNilBodyIsntChanged(t *testing.T) {
	defer internalconfig.ResetConfig()

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Nil(t, r.Body)
	})

	wh, _ := WrapHandler(h, mock.SpanFromContext, &Options{}, map[string]string{}).(*handler)
	wh.dataCaptureConfig = emptyTestConfig

	ih := &mockHandler{baseHandler: wh}

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", nil)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)

	assert.Equal(t, 1, len(ih.spans))
}

func TestServerRequestIsSuccessfullyTraced(t *testing.T) {
	defer internalconfig.ResetConfig()

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("request_id", "abc123xyz")
		rw.WriteHeader(202)
		rw.Write([]byte("test_res"))
		rw.Write([]byte("ponse_body"))
	})

	wh, _ := WrapHandler(h, mock.SpanFromContext, &Options{}, map[string]string{"foo": "bar"}).(*handler)
	wh.dataCaptureConfig = emptyTestConfig
	ih := &mockHandler{baseHandler: wh}

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("test_request_body"))
	r.Header.Add("api_key", "xyz123abc")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	assert.Equal(t, "test_response_body", w.Body.String())

	assert.Equal(t, 1, len(ih.spans))

	span := ih.spans[0]
	assert.Equal(t, "http://traceable.ai/foo?user_id=1", span.ReadAttribute("http.url").(string))
	assert.Equal(t, "traceable.ai", span.ReadAttribute("http.request.header.host"))
	assert.Equal(t, "bar", span.ReadAttribute("foo"))

	_ = span.ReadAttribute("container_id") // needed in containarized envs
	assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
}

func TestHostIsSuccessfullyRecorded(t *testing.T) {
	defer internalconfig.ResetConfig()

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Nil(t, r.Body)
	})

	wh, _ := WrapHandler(h, mock.SpanFromContext, &Options{}, map[string]string{}).(*handler)
	wh.dataCaptureConfig = emptyTestConfig

	ih := &mockHandler{baseHandler: wh}

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", nil)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	span := ih.spans[0]
	assert.NotNil(t, span)

	hostHeaderValue, ok := span.Attributes["http.request.header.host"]
	assert.True(t, ok)
	assert.Equal(t, hostHeaderValue, "traceable.ai")
}

func TestServerRequestHeadersAreSuccessfullyRecorded(t *testing.T) {
	defer internalconfig.ResetConfig()

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

		wh, _ := WrapHandler(h, mock.SpanFromContext, &Options{}, map[string]string{}).(*handler)
		ih := &mockHandler{baseHandler: wh}
		wh.dataCaptureConfig = emptyTestConfig
		wh.dataCaptureConfig.HttpHeaders = &config.Message{
			Request:  config.Bool(tCase.captureHTTPHeadersRequestConfig),
			Response: config.Bool(tCase.captureHTTPHeadersResponseConfig),
		}

		r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("test_request_body"))
		r.Header.Add("api_key", "xyz123abc")
		w := httptest.NewRecorder()

		ih.ServeHTTP(w, r)

		spans := ih.spans
		assert.Equal(t, 1, len(spans))

		span := spans[0]
		if tCase.captureHTTPHeadersRequestConfig {
			assert.Equal(t, "xyz123abc", span.ReadAttribute("http.request.header.api_key").(string))
		} else {
			assert.Nil(t, span.ReadAttribute("http.request.header.Api_key"))
		}

		if tCase.captureHTTPHeadersResponseConfig {
			assert.Equal(t, "abc123xyz", span.ReadAttribute("http.response.header.request_id").(string))
		} else {
			assert.Nil(t, span.ReadAttribute("http.response.header.request_id"))
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
		shouldBase64EncodeRequestBody  bool
		shouldBase64EncodeResponseBody bool
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
			responseContentType:            "application/json; charset=utf-8",
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
		"request has multipart/form-data content type and body with config enabled": {
			captureHTTPBodyConfig:          true,
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "multipart/form-data",
			responseContentType:            "application/json; charset=utf-8",
			shouldHaveRecordedRequestBody:  true,
			shouldHaveRecordedResponseBody: true,
			shouldBase64EncodeRequestBody:  true,
			shouldBase64EncodeResponseBody: false,
		},
		"response has multipart/form-data content type and body with config enabled": {
			captureHTTPBodyConfig:          true,
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "application/json; charset=utf-8",
			responseContentType:            "multipart/form-data",
			shouldHaveRecordedRequestBody:  true,
			shouldHaveRecordedResponseBody: true,
			shouldBase64EncodeRequestBody:  false,
			shouldBase64EncodeResponseBody: true,
		},
		"both request and response has multipart/form-data content type and body with config enabled": {
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
		"both request and response has multipart/form-data content type and body capture disabled": {
			captureHTTPBodyConfig:          false,
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "multipart/form-data",
			responseContentType:            "multipart/form-data",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
			shouldBase64EncodeRequestBody:  true,
			shouldBase64EncodeResponseBody: true,
		},
		"both request and response has multipart/form-data content type and multipart form-data body capture disabled": {
			captureHTTPBodyConfig:          true,
			requestBody:                    "test_request_body",
			responseBody:                   "test_response_body",
			requestContentType:             "multipart/form-data",
			responseContentType:            "multipart/form-data",
			shouldHaveRecordedRequestBody:  false,
			shouldHaveRecordedResponseBody: false,
			shouldBase64EncodeRequestBody:  false,
			shouldBase64EncodeResponseBody: false,
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Add("Content-Type", tCase.responseContentType)
				rw.Header().Add("Content-Type", "charset=UTF-8")
				rw.Write([]byte(tCase.responseBody))
			})

			wh, _ := WrapHandler(h, mock.SpanFromContext, &Options{}, map[string]string{}).(*handler)
			wh.dataCaptureConfig = emptyTestConfig
			wh.dataCaptureConfig.HttpBody = &config.Message{
				Request:  config.Bool(tCase.captureHTTPBodyConfig),
				Response: config.Bool(tCase.captureHTTPBodyConfig),
			}
			defaultAllowedContentTypes := internalconfig.GetConfig().DataCapture.AllowedContentTypes
			// add multipart/form-data to the allowed content types.
			if tCase.shouldBase64EncodeRequestBody || tCase.shouldBase64EncodeResponseBody {
				internalconfig.GetConfig().DataCapture.AllowedContentTypes = append(internalconfig.GetConfig().DataCapture.AllowedContentTypes,
					config.String("multipart/form-data"))
			}
			ih := &mockHandler{baseHandler: wh}

			r, _ := http.NewRequest("GET", "http://traceable.ai/foo", strings.NewReader(tCase.requestBody))
			r.Header.Add("content-type", tCase.requestContentType)

			w := httptest.NewRecorder()

			ih.ServeHTTP(w, r)

			span := ih.spans[0]

			if tCase.shouldHaveRecordedRequestBody {
				if tCase.shouldBase64EncodeRequestBody {
					assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte(tCase.requestBody)),
						span.ReadAttribute("http.request.body.base64").(string))
					assert.Nil(t, span.ReadAttribute("http.request.body"))
				} else {
					assert.Equal(t, tCase.requestBody, span.ReadAttribute("http.request.body").(string))
					assert.Nil(t, span.ReadAttribute("http.request.body.base64"))
				}
			} else {
				assert.Nil(t, span.ReadAttribute("http.request.body"))
				assert.Nil(t, span.ReadAttribute("http.request.body.base64"))
			}

			if tCase.shouldHaveRecordedResponseBody {
				if tCase.shouldBase64EncodeResponseBody {
					assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte(tCase.responseBody)),
						span.ReadAttribute("http.response.body.base64").(string))
					assert.Nil(t, span.ReadAttribute("http.response.body"))
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

func TestServerRequestFilter(t *testing.T) {
	tCases := map[string]struct {
		url                    string
		headerKeys             []string
		headerValues           []string
		body                   string
		options                *Options
		blocked                bool
		allowMultipartFormData bool
	}{
		"no filters": {
			options: &Options{},
			blocked: false,
		},
		"all filters no match, verify filter arguments": {
			url:          "http://localhost/foo",
			headerKeys:   []string{"content-type"},
			headerValues: []string{"application/json"},
			body:         "haha",
			options: &Options{
				Filter: mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
						assert.Equal(t, "http://localhost/foo", url)
						assert.Equal(t, 1, len(headers))
						assert.Equal(t, []string{"application/json"}, headers["Content-Type"])
						return result.FilterResult{}
					},
					BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult {
						assert.Equal(t, []byte("haha"), body)
						return result.FilterResult{}
					},
				},
			},
		},
		"all filters no match, verify filter arguments for multipart/form-data": {
			url:          "http://localhost/foo",
			headerKeys:   []string{"content-type"},
			headerValues: []string{"multipart/form-data"},
			body:         "haha",
			options: &Options{
				Filter: mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
						assert.Equal(t, "http://localhost/foo", url)
						assert.Equal(t, 1, len(headers))
						assert.Equal(t, []string{"multipart/form-data"}, headers["Content-Type"])
						return result.FilterResult{}
					},
					BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult {
						assert.Equal(t, []byte(base64.RawStdEncoding.EncodeToString([]byte("haha"))), body)
						return result.FilterResult{}
					},
				},
			},
			allowMultipartFormData: true,
		},
		"all filters no match, multipart body not captured, verify filter arguments for multipart/form-data": {
			url:          "http://localhost/foo",
			headerKeys:   []string{"content-type"},
			headerValues: []string{"multipart/form-data"},
			body:         "haha",
			options: &Options{
				Filter: mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
						assert.Equal(t, "http://localhost/foo", url)
						assert.Equal(t, 1, len(headers))
						assert.Equal(t, []string{"multipart/form-data"}, headers["Content-Type"])
						return result.FilterResult{}
					},
				},
			},
			allowMultipartFormData: false,
		},
		"url filter match": {
			url: "http://localhost/foo",
			options: &Options{
				Filter: mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
						return result.FilterResult{Block: true, ResponseStatusCode: 403}
					},
				},
			},
			blocked: true,
		},
		"headers filters match": {
			url: "http://localhost/foo",
			options: &Options{
				Filter: mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
						return result.FilterResult{Block: true, ResponseStatusCode: 403}
					},
				},
			},
			blocked: true,
		},
		"body filters match": {
			url:          "http://localhost/foo",
			headerKeys:   []string{"content-type"},
			headerValues: []string{"application/json"},
			body:         "haha",
			options: &Options{
				Filter: mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult {
						return result.FilterResult{Block: true, ResponseStatusCode: 403}
					},
				},
			},
			blocked: true,
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			defaultAllowedContentTypes := internalconfig.GetConfig().DataCapture.AllowedContentTypes
			// add multipart/form-data to the allowed content types.
			if tCase.allowMultipartFormData {
				internalconfig.GetConfig().DataCapture.AllowedContentTypes = append(internalconfig.GetConfig().DataCapture.AllowedContentTypes,
					config.String("multipart/form-data"))
			}
			h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusOK)
			})

			wh, _ := WrapHandler(h, mock.SpanFromContext, tCase.options, map[string]string{}).(*handler)
			ih := &mockHandler{baseHandler: wh}
			r, _ := http.NewRequest("POST", tCase.url, strings.NewReader(tCase.body))
			for i := 0; i < len(tCase.headerKeys); i++ {
				r.Header.Add(tCase.headerKeys[i], tCase.headerValues[i])
			}

			w := httptest.NewRecorder()

			ih.ServeHTTP(w, r)
			if !tCase.blocked {
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			} else {
				assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
			}
			// reset allowed content types config
			internalconfig.GetConfig().DataCapture.AllowedContentTypes = defaultAllowedContentTypes
		})
	}
}

func TestProcessingBodyIsTrimmed(t *testing.T) {
	defer internalconfig.ResetConfig()

	bodyMaxProcessingSizeBytes := 1

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})

	wh, _ := WrapHandler(h, mock.SpanFromContext, &Options{
		Filter: mock.Filter{
			BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult {
				assert.Equal(t, "{", string(body)) // body is truncated
				return result.FilterResult{Block: true, ResponseStatusCode: 403}
			},
		},
	}, map[string]string{}).(*handler)
	wh.dataCaptureConfig = emptyTestConfig
	wh.dataCaptureConfig.HttpBody.Request = config.Bool(true)
	wh.dataCaptureConfig.BodyMaxProcessingSizeBytes = config.Int32(int32(bodyMaxProcessingSizeBytes))

	ih := &mockHandler{baseHandler: wh}

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("{\"foo\":\"bar\"}"))
	r.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)
}
