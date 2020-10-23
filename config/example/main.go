package main

import (
	"fmt"
	"os"

	"github.com/traceableai/goagent/config"
)

func main() {
	os.Setenv("HT_CONFIG_FILE", "./example/traceable.json")
	os.Setenv("HT_DATA_CAPTURE_ENABLE_HTTP_PAYLOADS", "1")
	cfg := config.Load()
	cfg.ServiceName = "example"
	fmt.Println(cfg.JSONString())
}
