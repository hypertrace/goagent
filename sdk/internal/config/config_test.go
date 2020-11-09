package config

import (
	"testing"

	"github.com/hypertrace/goagent/config"
)

func TestConfig(t *testing.T) {
	InitConfig(&config.AgentConfig{
		ServiceName: config.String("my_service"),
	})
}
