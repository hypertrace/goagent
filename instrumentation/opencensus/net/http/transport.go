package http

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opencensus"
	traceablehttp "github.com/hypertrace/goagent/sdk/net/http"
)

// WrapTransport returns a new http.RoundTripper that should be passed to
// the OpenCensus *ochttp.Transport
func WrapTransport(delegate http.RoundTripper) http.RoundTripper {
	return traceablehttp.WrapTransport(delegate, opencensus.SpanFromContext)
}
