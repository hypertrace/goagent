package config

import (
	"log"
	sync "sync"
)

var cfg *Config
var mux = sync.Mutex{}

// InitConfig initializes the config with default values
func InitConfig(c *Config) {
	if cfg != nil {
		log.Println("config already initialized, ignoring new config.")
		return
	}

	cfg = c
}

// GetConfig returns the config value
func GetConfig() *Config {
	if cfg == nil {
		mux.Lock()
		cfg = &Config{}
		mux.Unlock()
	}

	return cfg
}
