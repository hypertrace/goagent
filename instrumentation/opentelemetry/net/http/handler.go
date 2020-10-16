package http

import (
	"net/http"

	"github.com/traceableai/goagent/instrumentation/opentelemetry"
	traceablehttp "github.com/traceableai/goagent/sdk/net/http"
)

// WrapHandler returns a new round tripper instrumented that relies on the
// needs to be used with OTel instrumentation.
func WrapHandler(delegate http.Handler) http.Handler {
	return traceablehttp.WrapHandler(delegate, opentelemetry.SpanFromContext)
}
