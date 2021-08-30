package hypertrace // import "github.com/hypertrace/goagent/instrumentation/hypertrace"

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
)

type Span opentelemetry.Span

var (
	SpanFromContext = opentelemetry.SpanFromContext
	StartSpan       = opentelemetry.StartSpan
)
