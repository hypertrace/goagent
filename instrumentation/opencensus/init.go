package opencensus // import "github.com/hypertrace/goagent/instrumentation/opencensus"

import (
	"crypto/tls"
	"log"
	"net/http"

	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	config "github.com/hypertrace/agent-config/gen/go/v1"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	zipkin "github.com/openzipkin/zipkin-go"
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
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: !cfg.GetReporting().GetSecure().GetValue(),
		},
	}}

	reporter := zipkinHTTP.NewReporter(cfg.GetReporting().GetEndpoint().GetValue(), zipkinHTTP.Client(client))

	exporter := oczipkin.NewExporter(reporter, localEndpoint)

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return func() {
		err := reporter.Close()
		if err != nil {
			log.Printf("error while closing reporter: %v\n", err)
		}
	}
}
