package config

import (
	"github.com/traceableai/goagent/config"
	internalconfig "github.com/traceableai/goagent/sdk/internal/config"
)

// InitConfig allows users to initialize the config
func InitConfig(c config.AgentConfig) {
	internalconfig.InitConfig(c)
}
