package config

import (
	"log"
	sync "sync"
)

var cfg *AgentConfig
var mux = sync.Mutex{}

// String returns a *string value that can be used to set into the configuration
func String(s string) *string {
	return &s
}

// Bool return a *bool value that can be used to set into the configuration
func Bool(b bool) *bool {
	return &b
}

// InitConfig initializes the config with default values
func InitConfig(c *AgentConfig) {
	if cfg != nil {
		log.Println("config already initialized, ignoring new config.")
		return
	}

	mux.Lock()
	cfg = c
	mux.Lock()
}

// GetConfig returns the config value
func GetConfig() *AgentConfig {
	if cfg == nil {
		mux.Lock()
		cfg = &AgentConfig{}
		mux.Unlock()
	}

	return cfg
}
