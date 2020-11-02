package htiohttp

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opencensus"
	sdkhttp "github.com/hypertrace/goagent/sdk/net/http"
)

// WrapHandler returns a new http.Handler that should be passed to
// the *ochttp.Handler
func WrapHandler(delegate http.Handler) http.Handler {
	return sdkhttp.WrapHandler(delegate, opencensus.SpanFromContext)
}
