package config

import (
	"testing"

	"github.com/traceableai/goagent/config"
)

func TestConfig(t *testing.T) {
	InitConfig(config.AgentConfig{
		ServiceName: config.StringVal("my_service"),
	})
}
