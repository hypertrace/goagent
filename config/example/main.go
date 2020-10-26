package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/traceableai/goagent/config"
)

func main() {
	// defines the config file path
	os.Setenv("HT_CONFIG_FILE", "./example/config.json")

	// defines the DataCapture.EnableHTTPPayloads = true
	os.Setenv("HT_DATA_CAPTURE_ENABLE_HTTP_PAYLOADS", "false")

	// loads the config
	cfg := config.Load()

	// overrides statically the service name
	cfg.ServiceName = config.StringVal("example")

	// prints the config as a JSON
	c, err := json.Marshal(&cfg)
	if err != nil {
		log.Fatalf("failed to marshal config: %v", err)
	}

	fmt.Println(string(c))
}
