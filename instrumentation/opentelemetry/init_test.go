package opentelemetry

import "github.com/hypertrace/goagent/config"

func ExampleInit() {
	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.DataCapture.HttpHeaders.Request = config.Bool(true)
	cfg.Reporting.Address = config.String("api.traceable.ai")

	shutdown := Init(cfg)
	defer shutdown()
}
