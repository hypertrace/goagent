package config

// LoadFromEnv loads env and default values on existing AgentConfig instance
func (x *AgentConfig) LoadFromEnv() {
	x.loadFromEnv("HT_", &defaultConfig)
}
