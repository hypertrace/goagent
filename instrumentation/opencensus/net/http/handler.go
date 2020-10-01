package http

import (
	"net/http"

	"github.com/traceableai/goagent/instrumentation/opencensus"
	traceablehttp "github.com/traceableai/goagent/sdk/net/http"
)

// EnrichHandler returns a new round tripper instrumented that relies on the
// needs to be used with OTel instrumentation.
func EnrichHandler(delegate http.Handler) http.Handler {
	return traceablehttp.EnrichHandler(delegate, opencensus.SpanFromContext)
}
