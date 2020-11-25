package hypermux

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkhttp "github.com/hypertrace/goagent/sdk/net/http"
	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
	"go.opentelemetry.io/otel/api/global"
)

func spanNameFormatter(operation string, r *http.Request) (spanName string) {
	route := mux.CurrentRoute(r)
	if route != nil {
		var err error
		spanName, err = route.GetPathTemplate()
		if err != nil {
			spanName, _ = route.GetPathRegexp()
		}
	}

	if spanName == "" {
		// if somehow retrieving the path template or path regexp fails, we still
		// want to use the method as fallback.
		spanName = r.Method
	}

	return
}

// NewMiddleware sets up a handler to start tracing the incoming requests.
func NewMiddleware() mux.MiddlewareFunc {
	return func(delegate http.Handler) http.Handler {
		return otelhttp.NewHandler(
			sdkhttp.WrapHandler(delegate, opentelemetry.SpanFromContext),
			"",
			otelhttp.WithTracer(global.Tracer(opentelemetry.TracerDomain)),
			otelhttp.WithSpanNameFormatter(spanNameFormatter),
		)
	}
}
