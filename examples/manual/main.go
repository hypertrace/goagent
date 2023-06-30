//go:build ignore

package main

import (
	"context"
	"time"

	agentconfig "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/sdk"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("manual-spans")
	cfg.Reporting.Endpoint = agentconfig.String("localhost:4317")
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP

	flusher := hypertrace.Init(cfg)
	defer flusher()

	_, span, ender := hypertrace.StartSpan(context.Background(), "manual-span", &sdk.SpanOptions{Kind: sdk.SpanKindServer})
	span.SetAttribute("manual-attr-0", "attr-val-0")
	span.SetAttribute("manual-attr-1", 123)

	span.AddEvent("event0", time.Now(), map[string]interface{}{})
	span.AddEvent("event1", time.Now(), map[string]interface{}{"k1": "v1", "k2": 3.142})

	span.SetStatus(sdk.StatusCodeOk, "")
	ender()

	time.Sleep(5 * time.Second)
}
