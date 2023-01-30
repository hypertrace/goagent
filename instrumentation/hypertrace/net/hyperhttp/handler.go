package hyperhttp // import "github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewHandler wraps the passed handler, functioning like middleware.
func NewHandler(base http.Handler, operation string, opts ...Option) http.Handler {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return otelhttp.NewHandler(
		sdkhttp.WrapHandler(base, operation, opentelemetry.SpanFromContext, o.toSDKOptions(), map[string]string{}),
		operation,
	)
}
