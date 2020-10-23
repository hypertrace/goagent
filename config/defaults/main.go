package main

import (
	"fmt"

	"github.com/traceableai/goagent/config"
)

func main() {
	cfg := config.Load()
	fmt.Println(cfg.JSONString())
}
