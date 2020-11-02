package examples

import (
	"github.com/traceableai/goagent/config"
	"github.com/traceableai/goagent/instrumentation/opentelemetry"
)

// InitTracer initializes the tracer and register it globally
func InitTracer(serviceName string) func() {
	cfg := config.Load()

	cfg.ServiceName = config.StringVal(serviceName)

	cfg.DataCapture.HTTPHeaders.Request = config.BoolVal(true)
	cfg.DataCapture.HTTPHeaders.Response = config.BoolVal(true)
	cfg.DataCapture.HTTPBody.Request = config.BoolVal(true)
	cfg.DataCapture.HTTPBody.Response = config.BoolVal(true)

	cfg.Reporting.Address = config.StringVal("localhost")
	cfg.Reporting.Secure = config.BoolVal(false)

	return opentelemetry.Init(cfg)
}
