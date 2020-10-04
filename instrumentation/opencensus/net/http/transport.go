package http

import (
	"net/http"

	"github.com/traceableai/goagent/instrumentation/opencensus"
	traceablehttp "github.com/traceableai/goagent/sdk/net/http"
)

// EnrichTransport returns a new round tripper instrumented that relies on the
// needs to be used with OpenCensus instrumentation.
func EnrichTransport(delegate http.RoundTripper) http.RoundTripper {
	return traceablehttp.EnrichTransport(delegate, opencensus.SpanFromContext)
}
