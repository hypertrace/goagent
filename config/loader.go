package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const defaultConfigFile = "./config.json"

func getBoolEnv(s string) (bool, bool) {
	if val := os.Getenv(s); val != "" {
		return val == "true", true
	}

	return false, false
}

func getStringEnv(s string) (string, bool) {
	if val := os.Getenv(s); val != "" {
		return val, true
	}

	return "", false
}

func loadFromFile(c *AgentConfig, filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, c)
}

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
