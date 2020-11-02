package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const defaultConfigFile = "./config.json"

// getBoolEnv returns the bool value for an env var and a confirmation
// if the var exists
func getBoolEnv(name string) (bool, bool) {
	if val := os.Getenv(name); val != "" {
		return val == "true", true
	}

	return false, false
}

// getStringEnv returns the string value for an env var and a confirmation
// if the var exists
func getStringEnv(name string) (string, bool) {
	if val := os.Getenv(name); val != "" {
		return val, true
	}

	return "", false
}

// loadFromFile loads the agent config from a file
func loadFromFile(c *AgentConfig, filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, c)
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Load loads the configuration from the default values, config file and env vars.
func Load() AgentConfig {
	cfg := AgentConfig{}

	if configFile := os.Getenv("HT_CONFIG_FILE"); configFile == "" {
		if fileExists(defaultConfigFile) {
			loadFromFile(&cfg, defaultConfigFile)
		}
	} else {
		absConfigFile, err := filepath.Abs(configFile)
		if err != nil {
			log.Printf("failed to resolve absolute path for %q: %v.\n", configFile, err)
		}

		if fileExists(absConfigFile) {
			loadFromFile(&cfg, absConfigFile)
		} else {
			log.Printf("config file %q does not exist.\n", absConfigFile)
		}
	}

	cfg.loadFromEnv("HT_", &defaultConfig)

	return cfg
}
