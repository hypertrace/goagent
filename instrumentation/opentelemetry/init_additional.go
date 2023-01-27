package opentelemetry

import (
	"context"
	"log"
	"strings"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	modbsp "github.com/hypertrace/goagent/instrumentation/opentelemetry/batchspanprocessor"
	sdkconfig "github.com/hypertrace/goagent/sdk/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// InitAsAdditional initializes opentelemetry tracing and returns a span processor and a shutdown
// function to flush data immediately on a termination signal.
// This is ideal for when we use goagent along with other opentelemetry setups.
func InitAsAdditional(cfg *config.AgentConfig) (trace.SpanProcessor, func()) {
	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return nil, func() {}
	}
	sdkconfig.InitConfig(cfg)

	exporterFactory = makeExporterFactory(cfg)

	exporter, err := exporterFactory()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.GetServiceName().GetValue() != "" {
		resource, err := resource.New(
			context.Background(),
			resource.WithAttributes(createResources(cfg.GetServiceName().GetValue(), cfg.ResourceAttributes)...),
		)
		if err != nil {
			log.Fatal(err)
		}

		exporter = addResourceToSpans(exporter, resource)
	}

	return modbsp.CreateBatchSpanProcessor(
			cfg.GetTelemetry() != nil && cfg.GetTelemetry().GetMetricsEnabled().GetValue(), // metrics enabled
			exporter,
			trace.WithBatchTimeout(batchTimeout)),
		func() {
			err := exporter.Shutdown(context.Background())
			if err != nil {
				log.Printf("error while shutting down exporter: %v\n", err)
			}
			sdkconfig.ResetConfig()
		}
}

type shieldResourceSpan struct {
	trace.ReadOnlySpan
	resource *resource.Resource
}

var _ trace.ReadOnlySpan = (*shieldResourceSpan)(nil)

func (s *shieldResourceSpan) Resource() *resource.Resource {
	return s.resource
}

type resourcePutter struct {
	trace.SpanExporter
	resource *resource.Resource
}

func (e *resourcePutter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	newSpans := []trace.ReadOnlySpan{}
	for _, span := range spans {
		newSpans = append(newSpans, &shieldResourceSpan{span, e.resource})
	}

	return e.SpanExporter.ExportSpans(ctx, newSpans)
}

func addResourceToSpans(e trace.SpanExporter, r *resource.Resource) trace.SpanExporter {
	return &resourcePutter{SpanExporter: e, resource: r}
}

// shieldAttrsSpan is a wrapper around a span that removes all attributes added by goagent as not
// all the the tracing servers can handle the load of big attributes like body or headers.
type shieldAttrsSpan struct {
	trace.ReadOnlySpan
	prefixes []string
}

func (s *shieldAttrsSpan) Attributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{}
	for _, attr := range s.ReadOnlySpan.Attributes() {
		key := string(attr.Key)
		hasPrefix := false
		for _, prefix := range s.prefixes {
			if strings.HasPrefix(key, prefix) {
				hasPrefix = true
				break
			}
		}

		if !hasPrefix {
			attrs = append(attrs, attr)
		}
	}

	return attrs
}

type attrsRemover struct {
	trace.SpanExporter
	prefixes []string
}

func (e *attrsRemover) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	newSpans := []trace.ReadOnlySpan{}
	for _, span := range spans {
		newSpans = append(newSpans, &shieldAttrsSpan{span, e.prefixes})
	}

	return e.SpanExporter.ExportSpans(ctx, newSpans)
}

var attrsRemovalPrefixes = []string{
	"http.request.header.",
	"http.response.header.",
	"http.request.body",
	"http.response.body",
	"rpc.request.metadata.",
	"rpc.response.metadata.",
	"rpc.request.body",
	"rpc.response.body",
}

var RemoveGoAgentAttrs = MakeRemoveGoAgentAttrs(attrsRemovalPrefixes)

// RemoveGoAgentAttrs removes custom goagent attributes from the spans so that other tracing servers
// don't receive them and don't have to handle the load.
func MakeRemoveGoAgentAttrs(attrsRemovalPrefixes []string) func(sp trace.SpanExporter) trace.SpanExporter {
	return func(sp trace.SpanExporter) trace.SpanExporter {
		if sp == nil {
			return sp
		}

		return &attrsRemover{sp, attrsRemovalPrefixes}
	}
}
