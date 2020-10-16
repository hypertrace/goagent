package http

import (
	"net/http"

	"github.com/traceableai/goagent/instrumentation/opencensus"
	traceablehttp "github.com/traceableai/goagent/sdk/net/http"
)

// WrapHandler returns a new http.Handler that should be passed to
// the *ochttp.Handler
func WrapHandler(delegate http.Handler) http.Handler {
	return traceablehttp.WrapHandler(delegate, opencensus.SpanFromContext)
}
