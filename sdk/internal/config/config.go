package config

import (
	"log"
	"sync"

	"github.com/hypertrace/goagent/config"
	"google.golang.org/protobuf/proto"
)

var cfg *config.AgentConfig
var cfgMux = &sync.Mutex{}

// InitConfig initializes the config with default values
func InitConfig(c *config.AgentConfig) {
	cfgMux.Lock()
	defer cfgMux.Unlock()

	if cfg != nil {
		log.Println("config already initialized, ignoring new config.")
		return
	}

	// The reason why we clone the message instead of reusing the one passed by the user
	// is because user might decide to change values in runtime and that is undesirable
	// without a proper API.
	var ok bool
	cfg, ok = proto.Clone(c).(*config.AgentConfig)
	if !ok {
		log.Fatal("failed to initialize config.")
	}
}

// GetConfig returns the config value
func GetConfig() *config.AgentConfig {
	if cfg == nil {
		InitConfig(config.Load())
	}

	return cfg
}
