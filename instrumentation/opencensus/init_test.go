package opencensus

import "github.com/traceableai/goagent/config"

func ExampleInit() {
	shutdown := Init(config.AgentConfig{
		ServiceName: config.StringVal("my_example_svc"),
		DataCapture: &config.DataCapture{
			HTTPHeaders: &config.Message{
				Request: config.BoolVal(true),
			},
		},
		Reporting: &config.Reporting{
			Address: config.StringVal("api.traceable.ai"),
		},
	})

	defer shutdown()
}
