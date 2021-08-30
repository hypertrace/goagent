package config // import "github.com/hypertrace/goagent/sdk/config"

import (
	agentconfig "github.com/hypertrace/agent-config/gen/go/v1"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
)

// InitConfig allows users to initialize the config
func InitConfig(c *agentconfig.AgentConfig) {
	internalconfig.InitConfig(c)
}
