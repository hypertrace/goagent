package hyperhttp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/instrumentation/opencensus/internal"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func TestMain(m *testing.M) {
	sdkconfig.InitConfig(&config.AgentConfig{})
}

func TestServerRequestIsSuccessfullyTraced(t *testing.T) {
	flusher := internal.InitTracer()

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("request_id", "abc123xyz")
		rw.WriteHeader(202)
		rw.Write([]byte("test_res"))
		rw.Write([]byte("ponse_body"))
	})

	ih := &ochttp.Handler{Handler: WrapHandler(h, &sdkhttp.Options{})}

	r, _ := http.NewRequest("GET", "http://traceable.ai/foo?user_id=1", strings.NewReader("test_request_body"))
	r.Header.Add("api_key", "xyz123abc")
	w := httptest.NewRecorder()

	ih.ServeHTTP(w, r)
	assert.Equal(t, "test_response_body", w.Body.String())

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, "/foo", span.Name)
	assert.Equal(t, trace.SpanKindServer, span.SpanKind)

	assert.Equal(t, "http://traceable.ai/foo?user_id=1", span.Attributes["http.url"].(string))
	assert.Equal(t, "xyz123abc", span.Attributes["http.request.header.Api_key"].(string))
	assert.Equal(t, "abc123xyz", span.Attributes["http.response.header.Request_id"].(string))
}

func TestServerRecordsRequestAndResponseBodyAccordingly(t *testing.T) {
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
			h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Add("Content-Type", tCase.responseContentType)
				rw.Header().Add("Content-Type", "charset=UTF-8")
				rw.Write([]byte(tCase.responseBody))
			})

			ih := &ochttp.Handler{Handler: WrapHandler(h, &sdkhttp.Options{})}

			r, _ := http.NewRequest("GET", "http://traceable.ai/foo", strings.NewReader(tCase.requestBody))
			r.Header.Add("content-type", tCase.requestContentType)

			w := httptest.NewRecorder()

			ih.ServeHTTP(w, r)

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

func TestRequestExtractsIncomingHeadersSuccessfully(t *testing.T) {
	flusher := internal.InitTracer()

	h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})

	ih := &ochttp.Handler{Handler: WrapHandler(h, &sdkhttp.Options{})}

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
