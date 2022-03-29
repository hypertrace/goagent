package hyperhttp // import "github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewTransport wraps the provided http.RoundTripper with one that
// starts a span and injects the span context into the outbound request headers.
func NewTransport(base http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(
		sdkhttp.WrapTransport(base, opentelemetry.SpanFromContext, map[string]string{}),
	)
}
