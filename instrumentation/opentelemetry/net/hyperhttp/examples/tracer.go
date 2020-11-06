package examples

import (
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
)

// InitTracer initializes the tracer and register it globally
func InitTracer(serviceName string) func() {
	cfg := config.Load()

	cfg.ServiceName = serviceName

	cfg.Reporting.Address = config.StringVal("localhost")
	cfg.Reporting.Secure = config.BoolVal(false)

	return opentelemetry.Init(cfg)
}
