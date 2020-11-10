package hyperhttp

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkhttp "github.com/hypertrace/goagent/sdk/net/http"
	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
	"go.opentelemetry.io/otel/api/global"
)

const domain = "org.hypertrace.goagent"

// NewHandler wraps the passed handler, functioning like middleware.
func NewHandler(base http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(
		sdkhttp.WrapHandler(base, opentelemetry.SpanFromContext),
		operation,
		otelhttp.WithTracer(global.Tracer(domain)),
	)
}
