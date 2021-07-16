package hypergin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
	"go.opentelemetry.io/otel/trace"
	"gotest.tools/assert"
)

func handler(c *gin.Context) {
	c.Header("request_id", "xyz123abc")
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"id":   "123",
	})
}

func TestSpanRecordedCorrectly(t *testing.T) {
	_, flusher := internal.InitTracer()

	r := gin.Default()
	r.Use(Middleware(&sdkhttp.Options{}))
	r.POST("/things/:thing_id", handler)

	server := httptest.NewServer(r)
	defer server.Close()

	req := httptest.NewRequest(
		"POST",
		"http://example.com/things/123?include_something=1",
		bytes.NewBufferString(`{"name":"Jacinto"}`),
	)
	req.Header.Set("api_key", "abc123xyz")
	req.Header.Set("content-type", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if want, have := 200, w.Result().StatusCode; want != have {
		t.Errorf("unexpected status code, want %q, have %q", want, have)
	}

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, "/things/:thing_id", span.Name)
	assert.Equal(t, span.SpanKind, trace.SpanKindServer)

	attrs := internal.LookupAttributes(span.Attributes)
	assert.Equal(t, "POST", attrs.Get("http.method").AsString())
	assert.Equal(t, "abc123xyz", attrs.Get("http.request.header.api_key").AsString())
	assert.Equal(t, `{"name":"Jacinto"}`, attrs.Get("http.request.body").AsString())
	assert.Equal(t, "xyz123abc", attrs.Get("http.response.header.request_id").AsString())
	assert.Equal(t, `{"code":200,"id":"123"}`, attrs.Get("http.response.body").AsString())
	assert.Equal(t, "application/json; charset=utf-8", attrs.Get("http.response.header.content-type").AsString())
}

func TestTraceContextIsPropagated(t *testing.T) {
	_, flusher := internal.InitTracer()

	// Configure Gin server
	r := gin.Default()
	r.Use(Middleware(&sdkhttp.Options{}))
	r.POST("/things/:thing_id", handler)

	server := &http.Server{Addr: ":60543", Handler: r}
	defer server.Close()

	go server.ListenAndServe()

	// Configure http Client
	client := http.Client{
		Transport: hyperhttp.NewTransport(
			http.DefaultTransport,
		),
	}

	req, _ := http.NewRequest("POST",
		"http://localhost:60543/things/123",
		bytes.NewBufferString(`{"name":"Jacinto"}`))

	res, err := client.Do(req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if want, have := 200, res.StatusCode; want != have {
		t.Errorf("unexpected status code, want %q, have %q", want, have)
	}

	spans := flusher()
	if spans == nil {
		t.Errorf("failed")
	}

	assert.Equal(t, 1, len(spans))

	span := spans[0]
	attrs := internal.LookupAttributes(span.Attributes)
	b3RequestHeader := attrs.Get("http.request.header.b3").AsString()
	expectedHeader := span.SpanContext.TraceID().String() + "-" + span.ParentSpanID.String() + "-1"
	assert.Equal(t, b3RequestHeader, expectedHeader)

}
