package opencensus

import "github.com/hypertrace/goagent/config"

func ExampleInit() {
	cfg := config.Load()
	cfg.ServiceName = config.String("my_example_svc")
	cfg.DataCapture.HttpHeaders.Request = config.Bool(true)
	cfg.Reporting.Endpoint = config.String("https://api.traceable.ai:9411/api/v2/spans")

	shutdown := Init(cfg)
	defer shutdown()
}
