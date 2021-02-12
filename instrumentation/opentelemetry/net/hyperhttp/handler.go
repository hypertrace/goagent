package hyperhttp

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkhttp "github.com/hypertrace/goagent/sdk/net/http"
)

// WrapHandler returns a new round tripper instrumented that relies on the
// needs to be used with OTel instrumentation.
func WrapHandler(delegate http.Handler, options *sdkhttp.Options) http.Handler {
	return sdkhttp.WrapHandler(delegate, opentelemetry.SpanFromContext, options)
}
