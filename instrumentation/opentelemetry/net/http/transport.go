package http

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	traceablehttp "github.com/hypertrace/goagent/sdk/net/http"
)

// WrapTransport wraps an uninstrumented RoundTripper (e.g. http.DefaultTransport)
// and returns an instrumented RoundTripper that has to be used as base for the
// OTel's RoundTripper.
func WrapTransport(delegate http.RoundTripper) http.RoundTripper {
	return traceablehttp.WrapTransport(delegate, opentelemetry.SpanFromContext)
}
