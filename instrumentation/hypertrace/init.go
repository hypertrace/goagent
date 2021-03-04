package hypertrace

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
)

// Init initializes hypertrace tracing and returns a shutdown function to flush data immediately
// on a termination signal.
var Init = opentelemetry.Init

var InitWithResources = opentelemetry.InitWithResources