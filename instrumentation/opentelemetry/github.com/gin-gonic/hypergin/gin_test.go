package hypergin

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, "/things/:thing_id", span.Name())
	assert.Equal(t, span.SpanKind(), trace.SpanKindServer)

	attrs := internal.LookupAttributes(span.Attributes())
	assert.Equal(t, "POST", attrs.Get("http.method").AsString())
	assert.Equal(t, "abc123xyz", attrs.Get("http.request.header.api_key").AsString())
	assert.Equal(t, `{"name":"Jacinto"}`, attrs.Get("http.request.body").AsString())
	assert.Equal(t, "xyz123abc", attrs.Get("http.response.header.request_id").AsString())
	assert.Equal(t, `{"code":200,"id":"123"}`, attrs.Get("http.response.body").AsString())
	assert.Equal(t, "application/json; charset=utf-8", attrs.Get("http.response.header.content-type").AsString())
}

// Client -> GET Server1/send_thing_request -> POST Server2/things/:thing_id
func TestTraceContextIsPropagated(t *testing.T) {
	_, flusher := internal.InitTracer()

	var client = http.Client{
		Transport: hyperhttp.NewTransport(
			http.DefaultTransport,
		),
	}

	// Configure Gin server
	r := gin.Default()
	r.Use(Middleware(&sdkhttp.Options{}))
	r.POST("/things/:thing_id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"thing": "go",
		})
	})

	r2 := gin.Default()
	r2.Use(Middleware(&sdkhttp.Options{}))
	r2.GET("/send_thing_request", func(c *gin.Context) {
		req, _ := http.NewRequest("POST",
			"http://localhost:60543/things/123",
			bytes.NewBufferString(`{"name":"Jacinto"}`))

		req = req.WithContext(c.Request.Context())
		res, err := client.Do(req)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
			})
			return
		}
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		bodyString := string(bodyBytes)
		c.JSON(200, gin.H{
			"success":              true,
			"otherServiceResponse": bodyString,
		})
	})

	server := &http.Server{Addr: ":60543", Handler: r}
	defer server.Close()
	go server.ListenAndServe()

	server2 := &http.Server{Addr: ":60544", Handler: r2}
	defer server2.Close()
	go server2.ListenAndServe()

	req, _ := http.NewRequest("GET",
		"http://localhost:60544/send_thing_request", nil)

	res, err := client.Do(req)
	_, readErr := ioutil.ReadAll(res.Body)
	require.NoError(t, readErr)

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

	assert.Equal(t, 4, len(spans))
	assert.Equal(t, "/things/:thing_id", spans[0].Name())
	assert.Equal(t, spans[1].SpanContext().SpanID(), spans[0].Parent().SpanID())
	assert.Equal(t, "HTTP POST", spans[1].Name())
	assert.Equal(t, spans[2].SpanContext().SpanID(), spans[1].Parent().SpanID())
	assert.Equal(t, "/send_thing_request", spans[2].Name())
	assert.Equal(t, spans[3].SpanContext().SpanID(), spans[2].Parent().SpanID())
	assert.Equal(t, "HTTP GET", spans[3].Name())

	traceId := spans[0].SpanContext().TraceID().String()
	for _, span := range spans {
		assert.Equal(t, traceId, span.SpanContext().TraceID().String())
	}
}
