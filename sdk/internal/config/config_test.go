package config

import (
	"testing"

	config "github.com/hypertrace/agent-config/gen/go/v1"
)

func TestConfig(t *testing.T) {
	InitConfig(&config.AgentConfig{
		ServiceName: config.String("my_service"),
	})
}
