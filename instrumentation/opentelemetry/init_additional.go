package opentelemetry

import (
	"context"
	"log"
	"strings"

	config "github.com/hypertrace/agent-config/gen/go/v1"
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

	return trace.NewBatchSpanProcessor(exporter, trace.WithBatchTimeout(batchTimeout)), func() {
		exporter.Shutdown(context.Background())
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
}

func (s *shieldAttrsSpan) Attributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{}
	for _, attr := range s.ReadOnlySpan.Attributes() {
		if strings.HasPrefix(string(attr.Key), "http.request.header.") ||
			string(attr.Key) == "http.request.body" {
			continue
		}

		if strings.HasPrefix(string(attr.Key), "http.response.header.") ||
			string(attr.Key) == "http.response.body" {
			continue
		}

		if strings.HasPrefix(string(attr.Key), "rpc.request.metadata.") ||
			string(attr.Key) == "rpc.request.body" {
			continue
		}

		if strings.HasPrefix(string(attr.Key), "rpc.response.metadata.") ||
			string(attr.Key) == "rpc.response.body" {
			continue
		}

		attrs = append(attrs, attr)
	}

	return attrs
}

type attrsRemover struct {
	trace.SpanExporter
}

func (e *attrsRemover) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	newSpans := []trace.ReadOnlySpan{}
	for _, span := range spans {
		newSpans = append(newSpans, &shieldAttrsSpan{span})
	}

	return e.SpanExporter.ExportSpans(ctx, newSpans)
}

// RemoveGoAgentAttrs removes custom goagent attributes from the spans so that other tracing servers
// don't receive them and don't have to handle the load.
func RemoveGoAgentAttrs(sp trace.SpanExporter) trace.SpanExporter {
	if sp == nil {
		return sp
	}

	return &attrsRemover{sp}
}
