package opencensus

import (
	"fmt"

	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	"github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/traceableai/goagent/config"
	sdkconfig "github.com/traceableai/goagent/sdk/config"
	"go.opencensus.io/trace"
)

// Init initializes opencensus tracing and returns a shutdown function to flush data immediately
// on a termination signal.
func Init(cfg config.AgentConfig) func() {
	sdkconfig.InitConfig(cfg)
	localEndpoint, _ := zipkin.NewEndpoint(cfg.GetServiceName(), "localhost")

	reporterURI := fmt.Sprintf("http://%s:9411/api/v2/spans", cfg.Reporting.GetTracesEndpointHost())
	reporter := zipkinHTTP.NewReporter(reporterURI)

	exporter := oczipkin.NewExporter(reporter, localEndpoint)

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return func() {
		reporter.Close()
	}
}
