package hypermux

import (
	otelmux "github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/gorilla/hypermux"
)

// NewMiddleware sets up a handler to start tracing the incoming requests.
var NewMiddleware = otelmux.NewMiddleware
