package hyperhttp

import (
	"net/http"

	"github.com/hypertrace/goagent/instrumentation/opencensus"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
)

// WrapTransport returns a new http.RoundTripper that should be passed to
// the OpenCensus *ochttp.Transport
func WrapTransport(delegate http.RoundTripper) http.RoundTripper {
	return sdkhttp.WrapTransport(delegate, opencensus.SpanFromContext)
}
