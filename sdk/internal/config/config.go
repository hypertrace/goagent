package config

import (
	"log"
	"sync"

	agentconfig "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/config"
	"google.golang.org/protobuf/proto"
)

var cfg *agentconfig.AgentConfig
var cfgMux = &sync.Mutex{}

// InitConfig initializes the config with default values
func InitConfig(c *agentconfig.AgentConfig) {
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
	cfg, ok = proto.Clone(c).(*agentconfig.AgentConfig)
	if !ok {
		log.Fatal("failed to initialize config.")
	}
}

// GetConfig returns the config value
func GetConfig() *agentconfig.AgentConfig {
	if cfg == nil {
		InitConfig(config.Load())
	}

	return cfg
}

func ResetConfig() {
	cfgMux.Lock()
	defer cfgMux.Unlock()
	cfg = nil
}
