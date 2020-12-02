package opencensus

import (
	"fmt"

	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	"github.com/hypertrace/goagent/config"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	"github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

// Init initializes opencensus tracing and returns a shutdown function to flush data immediately
// on a termination signal.
func Init(cfg *config.AgentConfig) func() {
	sdkconfig.InitConfig(cfg)
	localEndpoint, _ := zipkin.NewEndpoint(cfg.GetServiceName().GetValue(), "localhost")

	protocol := "http"
	if cfg.GetReporting().GetSecure().GetValue() {
		protocol = "https"
	}
	reporterURL := fmt.Sprintf("%s://%s:9411/api/v2/spans", protocol, cfg.GetReporting().GetEndpoint().GetValue())
	reporter := zipkinHTTP.NewReporter(reporterURL)

	exporter := oczipkin.NewExporter(reporter, localEndpoint)

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return func() {
		reporter.Close()
	}
}
