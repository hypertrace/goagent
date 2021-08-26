package opencensus

import (
	"crypto/tls"
	"net/http"

	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	config "github.com/hypertrace/agent-config/gen/go/hypertrace/agent/config/v1"
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

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !cfg.GetReporting().GetSecure().GetValue()},
	}}

	reporter := zipkinHTTP.NewReporter(cfg.GetReporting().GetEndpoint().GetValue(), zipkinHTTP.Client(client))

	exporter := oczipkin.NewExporter(reporter, localEndpoint)

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return func() {
		reporter.Close()
	}
}
