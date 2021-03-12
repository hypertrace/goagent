package config

import (
	"github.com/hypertrace/goagent/config"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
)

// InitConfig allows users to initialize the config
func InitConfig(c *config.AgentConfig) {
	internalconfig.InitConfig(c)
}

func GetReportingEndpointConfig() string {
	return internalconfig.GetConfig().GetReporting().GetEndpoint().GetValue()
}

func GetReportingSecureConfig() bool {
	return internalconfig.GetConfig().GetReporting().GetSecure().GetValue()
}
