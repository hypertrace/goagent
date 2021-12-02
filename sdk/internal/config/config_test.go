package config

import (
	"testing"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	InitConfig(&config.AgentConfig{
		ServiceName: config.String("my_service"),
	})
	defer ResetConfig()

	assert.Equal(t, "my_service", GetConfig().ServiceName.Value)
}
