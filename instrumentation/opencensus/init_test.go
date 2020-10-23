package opencensus

import "github.com/traceableai/goagent/config"

func ExampleInit() {
	shutdown := Init(config.AgentConfig{
		ServiceName: "my_example_svc",
		DataCapture: &config.DataCapture{
			EnableHTTPHeaders: true,
		},
		Reporting: &config.Reporting{
			TracesEndpointHost: "api.traceable.ai",
		},
	})

	defer shutdown()
}
