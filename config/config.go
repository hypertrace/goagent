package config

import (
	agentconfig "github.com/hypertrace/agent-config/gen/go/hypertrace/agent/config/v1"
)

func Load() *agentconfig.AgentConfig {
	return agentconfig.Load(
		agentconfig.WithDefaults(&defaultConfig),
	)
}

func LoadFromFile(configFile string) *agentconfig.AgentConfig {
	return agentconfig.LoadFromFile(
		configFile,
		agentconfig.WithDefaults(&defaultConfig),
	)
}

func LoadEnv(cfg *agentconfig.AgentConfig) {
	cfg.LoadFromEnv(agentconfig.WithDefaults(&defaultConfig))
}

func PropagationFormats(formats ...agentconfig.PropagationFormat) []agentconfig.PropagationFormat {
	return formats
}

var (
	Bool                           = agentconfig.Bool
	String                         = agentconfig.String
	Int32                          = agentconfig.Int32
	TraceReporterType_OTLP         = agentconfig.TraceReporterType_OTLP
	TraceReporterType_ZIPKIN       = agentconfig.TraceReporterType_ZIPKIN
	PropagationFormat_B3           = agentconfig.PropagationFormat_B3
	PropagationFormat_TRACECONTEXT = agentconfig.PropagationFormat_TRACECONTEXT
)
