package hypermux

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal"
	sdkhttp "github.com/hypertrace/goagent/sdk/net/http"
	"go.opentelemetry.io/otel/trace"
	"gotest.tools/assert"
)

// inspired in https://github.com/jcchavezs/httptest-php/blob/e6a65c73/src/HttpTest/HttpTestServer.php#L150
func findAvailablePort() (int, error) {
	for port := 60000; port < 65535; port++ {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		defer l.Close()

		if err == nil {
			return port, nil
		}
	}

	return 0, errors.New("failed to find an available port")
}

func TestSpanRecordedCorrectly(t *testing.T) {
	_, flusher := internal.InitTracer()

	r := mux.NewRouter()
	r.HandleFunc("/things/{thing_id}", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("content-type", "application/json")
		rw.Header().Add("content-type", "charset=utf-8")
		rw.Header().Set("request_id", "xyz123abc")
		rw.WriteHeader(202)
		rw.Write([]byte(`{"id":123}`))
	})
	r.Use(NewMiddleware(&sdkhttp.Options{}))

	port, err := findAvailablePort()
	if err != nil {
		t.Fatal(err)
	}

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r}
	defer server.Close()

	go server.ListenAndServe()

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:%d/things/123?include_something=1", port),
		bytes.NewBufferString(`{"name":"Jacinto"}`),
	)
	req.Header.Set("api_key", "abc123xyz")
	req.Header.Set("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if want, have := 202, res.StatusCode; want != have {
		t.Errorf("unexpected status code, want %q, have %q", want, have)
	}

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, "/things/{thing_id}", span.Name)
	assert.Equal(t, span.SpanKind, trace.SpanKindServer)

	attrs := internal.LookupAttributes(span.Attributes)
	assert.Equal(t, "POST", attrs.Get("http.method").AsString())
	assert.Equal(t, "abc123xyz", attrs.Get("http.request.header.api_key").AsString())
	assert.Equal(t, `{"name":"Jacinto"}`, attrs.Get("http.request.body").AsString())
	assert.Equal(t, "xyz123abc", attrs.Get("http.response.header.request_id").AsString())
	assert.Equal(t, `{"id":123}`, attrs.Get("http.response.body").AsString())
	assert.Equal(t, "application/json", attrs.Get("http.response.header.content-type[0]").AsString())
	assert.Equal(t, "charset=utf-8", attrs.Get("http.response.header.content-type[1]").AsString())
}
