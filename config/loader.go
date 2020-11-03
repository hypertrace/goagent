package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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

	switch ext := filepath.Ext(filename); ext {
	case ".json":
		return json.Unmarshal(content, c)
	case ".yaml", ".yml":
		return yaml.Unmarshal(content, c)
	default:
		return fmt.Errorf("unknown extension: %s", ext)
	}
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

	if configFile := os.Getenv("HT_CONFIG_FILE"); configFile != "" {
		absConfigFile, err := filepath.Abs(configFile)
		if err != nil {
			log.Printf("failed to resolve absolute path for %q: %v.\n", configFile, err)
		}

		if fileExists(absConfigFile) {
			if err := loadFromFile(&cfg, absConfigFile); err != nil {
				log.Printf("failed to load the config fromt %q: %v\n", absConfigFile, err)
			}
		} else {
			log.Printf("config file %q does not exist.\n", absConfigFile)
		}
	}

	cfg.loadFromEnv("HT_", &defaultConfig)

	return cfg
}

// LoadFromFile loads the configuration from the default values, config file and env vars.
func LoadFromFile(configFile string) AgentConfig {
	cfg := AgentConfig{}

	absConfigFile, err := filepath.Abs(configFile)
	if err != nil {
		log.Printf("failed to resolve absolute path for %q: %v.\n", configFile, err)
	}

	if fileExists(absConfigFile) {
		if err := loadFromFile(&cfg, absConfigFile); err != nil {
			log.Printf("failed to load the config fromt %q: %v\n", absConfigFile, err)
		}
	} else {
		log.Printf("config file %q does not exist.\n", absConfigFile)
	}

	cfg.loadFromEnv("HT_", &defaultConfig)

	return cfg
}
