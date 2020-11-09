package examples

import (
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
)

// InitTracer initializes the tracer and register it globally
func InitTracer(serviceName string) func() {
	cfg := config.Load()

	cfg.ServiceName = config.String(serviceName)

	cfg.Reporting.Address = config.String("localhost")
	cfg.Reporting.Secure = config.Bool(false)

	return opentelemetry.Init(cfg)
}
