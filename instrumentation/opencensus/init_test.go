package opencensus

import "github.com/hypertrace/goagent/config"

func ExampleInit() {
	cfg := config.Load()
	cfg.ServiceName = config.StringVal("my_example_svc")
	cfg.DataCapture.HttpHeaders.Request = config.BoolVal(true)
	cfg.Reporting.Address = config.StringVal("api.traceable.ai")

	shutdown := Init(cfg)
	defer shutdown()
}
