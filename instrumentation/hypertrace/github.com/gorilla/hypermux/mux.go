package hypermux // import "github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/gorilla/hypermux"

import (
	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/gorilla/hypermux"
)

func NewMiddleware(opts ...Option) mux.MiddlewareFunc {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return hypermux.NewMiddleware(o.toSDKOptions())
}
