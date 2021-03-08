package config

// LoadFromEnv loads env and default values on existing AgentConfig instance
// where defaults only overrides empty values while env vars can override all
// of them.
func (x *AgentConfig) LoadFromEnv() {
	x.loadFromEnv("HT_", &defaultConfig)
}
