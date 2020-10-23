package config

import (
	"log"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/traceableai/goagent/config"
)

var cfg *config.AgentConfig
var mux = sync.Mutex{}

// InitConfig initializes the config with default values
func InitConfig(c config.AgentConfig) {
	if cfg != nil {
		log.Println("config already initialized, ignoring new config.")
		return
	}

	cfg = &config.AgentConfig{}
	mux.Lock()
	err := copier.Copy(cfg, &c)
	mux.Unlock()
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}
}

// GetConfig returns the config value
func GetConfig() *config.AgentConfig {
	if cfg == nil {
		InitConfig(config.Load())
	}

	return cfg
}
