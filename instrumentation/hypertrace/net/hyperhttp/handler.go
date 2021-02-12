package hyperhttp

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkhttp "github.com/hypertrace/goagent/sdk/net/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewHandler wraps the passed handler, functioning like middleware.
func NewHandler(base http.Handler, operation string, options *sdkhttp.Options) http.Handler {
	return otelhttp.NewHandler(
		sdkhttp.WrapHandler(base, opentelemetry.SpanFromContext, options),
		operation,
	)
}
