package hyperhttp // import "github.com/hypertrace/goagent/instrumentation/opentelemetry/net/hyperhttp"

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
)

// WrapTransport wraps an uninstrumented RoundTripper (e.g. http.DefaultTransport)
// and returns an instrumented RoundTripper that has to be used as base for the
// OTel's RoundTripper.
func WrapTransport(delegate http.RoundTripper) http.RoundTripper {
	return sdkhttp.WrapTransport(delegate, opentelemetry.SpanFromContext, map[string]string{})
}
