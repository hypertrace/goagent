package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
)

// getBoolEnv returns the bool value for an env var and a confirmation
// if the var exists
func getBoolEnv(name string) (bool, bool) {
	val := os.Getenv(name)
	switch val {
	case "true":
		return true, true
	case "false":
		return false, true
	default:
		return false, false
	}
}

// getStringEnv returns the string value for an env var and a confirmation
// if the var exists
func getStringEnv(name string) (string, bool) {
	if val := os.Getenv(name); val != "" {
		return val, true
	}

	return "", false
}

// getInt32Env returns the int32 value for an env var and a confirmation
// if the var exists
func getInt32Env(name string) (int32, bool) {
	if val := os.Getenv(name); val != "" {
		intVal, err := strconv.Atoi(val)
		return int32(intVal), err == nil
	}

	return 0, false
}

// loadFromFile loads the agent config from a file
func loadFromFile(c *AgentConfig, filename string) error {
	switch ext := filepath.Ext(filename); ext {
	case ".json":
		freader, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open file %q: %v", filename, err)
		}
		// The usage of wrappers for scalars in protos make it impossible to use standard
		// unmarshalers as the wrapped values aren't scalars but of type Message, hence they
		// have object structure in json e.g. myBoolVal: {Value: true} instead of myBoolVal:true
		// jsonpb is meant to solve this problem.
		return jsonpb.Unmarshal(freader, c)
	case ".yaml", ".yml":
		fcontent, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %v", filename, err)
		}
		// Because of the reson mentioned above we can't use YAML parsers either and hence
		// we convert the YAML into JSON in order to parse the JSON value with jsonpb.
		// The implications of this is that comments and multi-line strings aren't desirable.
		fcontentAsJSON, err := yaml.YAMLToJSON(fcontent)
		if err != nil {
			return fmt.Errorf("failed to parse file %q: %v", filename, err)
		}
		return jsonpb.UnmarshalString(string(fcontentAsJSON), c)
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
func Load() *AgentConfig {
	cfg := AgentConfig{}

	if configFile := os.Getenv("HT_CONFIG_FILE"); configFile != "" {
		absConfigFile, err := filepath.Abs(configFile)
		if err != nil {
			log.Printf("failed to resolve absolute path for %q: %v.\n", configFile, err)
		}

		if fileExists(absConfigFile) {
			if err := loadFromFile(&cfg, absConfigFile); err != nil {
				log.Printf("failed to load the config from %q: %v\n", absConfigFile, err)
			}
		} else {
			log.Printf("config file %q does not exist.\n", absConfigFile)
		}
	}

	cfg.loadFromEnv("HT_", &defaultConfig)

	return &cfg
}

// LoadFromFile loads the configuration from the default values, config file and env vars.
func LoadFromFile(configFile string) *AgentConfig {
	cfg := AgentConfig{}

	absConfigFile, err := filepath.Abs(configFile)
	if err != nil {
		log.Printf("failed to resolve absolute path for %q: %v.\n", configFile, err)
	}

	if fileExists(absConfigFile) {
		if err := loadFromFile(&cfg, absConfigFile); err != nil {
			log.Printf("failed to load the config from %q: %v\n", absConfigFile, err)
		}
	} else {
		log.Printf("config file %q does not exist.\n", absConfigFile)
	}

	cfg.loadFromEnv("HT_", &defaultConfig)

	return &cfg
}
