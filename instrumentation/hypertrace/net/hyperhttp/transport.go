package hyperhttp

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkhttp "github.com/hypertrace/goagent/sdk/net/http"
	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
	"go.opentelemetry.io/otel/api/global"
)

// NewTransport wraps the provided http.RoundTripper with one that
// starts a span and injects the span context into the outbound request headers.
func NewTransport(base http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(
		sdkhttp.WrapTransport(base, opentelemetry.SpanFromContext),
		otelhttp.WithTracer(global.Tracer(opentelemetry.TracerDomain)),
	)
}
