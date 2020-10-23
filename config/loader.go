package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const defaultConfig = "./config.json"

func getBoolEnv(s string, defaultValue bool) bool {
	if val := os.Getenv(s); val != "" {
		return val == "1"
	}

	return defaultValue
}

func getStringEnv(s string, defaultValue string) string {
	if val := os.Getenv(s); val != "" {
		return val
	}

	return defaultValue
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

// Load loads load
func Load() AgentConfig {
	cfg := AgentConfig{
		Reporting: &Reporting{
			TracesEndpointHost: getStringEnv("HT_REPORTING_TRACES_ENDPOINT_HOST", "localhost"),
		},
		DataCapture: &DataCapture{
			EnableHTTPPayloads: getBoolEnv("HT_DATA_CAPTURE_ENABLE_HTTP_PAYLOADS", false),
			EnableHTTPHeaders:  getBoolEnv("HT_DATA_CAPTURE_ENABLE_HTTP_HEADERS", false),
			EnableRPCPayloads:  getBoolEnv("HT_DATA_CAPTURE_ENABLE_HTTP_PAYLOADS", false),
			EnableRPCMetadata:  getBoolEnv("HT_DATA_CAPTURE_ENABLE_HTTP_METADATA", false),
		},
	}

	if configFile := os.Getenv("HT_CONFIG_FILE"); configFile == "" {
		if fileExists(defaultConfig) {
			loadFromFile(&cfg, defaultConfig)
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

	return cfg
}
