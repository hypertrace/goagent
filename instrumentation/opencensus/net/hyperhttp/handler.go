package hyperhttp // import "github.com/hypertrace/goagent/instrumentation/opencensus/net/hyperhttp"

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opencensus"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
)

// WrapHandler returns a new http.Handler that should be passed to
// the *ochttp.Handler
func WrapHandler(delegate http.Handler, options *sdkhttp.Options) http.Handler {
	// TODO: If I am doing this then I might have the metrics code in the wrong place.
	return sdkhttp.WrapHandler(delegate, "", opencensus.SpanFromContext, options, map[string]string{})
}
