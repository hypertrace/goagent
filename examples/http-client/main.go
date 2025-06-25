//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

type message struct {
	Content string `json:"message"`
}

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("client")
	cfg.Reporting.Endpoint = config.String("localhost:5442")
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP

	flusher := hypertrace.Init(cfg)
	defer flusher()

	client := http.Client{
		Transport: hyperhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequest("GET", "http://localhost:8081/foo", bytes.NewBufferString(`{"name":"こんにちは"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalf("failed to create the request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to perform the request: %v", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read the response body: %v", err)
	}
	defer res.Body.Close()

	m := &message{}
	err = json.Unmarshal(resBody, m)
	if err != nil {
		log.Fatalf("failed to unmarshal the response body: %v", err)
		return
	}

	fmt.Println(m.Content)
}
