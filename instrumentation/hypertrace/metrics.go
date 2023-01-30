package hypertrace // import "github.com/hypertrace/goagent/instrumentation/hypertrace"

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
)

var NewHttpOperationMetricsHandler = opentelemetry.NewHttpOperationMetricsHandler
