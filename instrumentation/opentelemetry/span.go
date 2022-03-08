package opentelemetry // import "github.com/hypertrace/goagent/instrumentation/opentelemetry"

import (
	"context"
	"fmt"
	"time"

	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ sdk.Span = (*Span)(nil)

type Span struct {
	trace.Span
}

func generateAttribute(key string, value interface{}) attribute.KeyValue {
	switch v := value.(type) {
	case bool:
		return attribute.Bool(key, v)
	case []bool:
		return attribute.BoolSlice(key, v)
	case int:
		return attribute.Int(key, v)
	case []int:
		return attribute.IntSlice(key, v)
	case int64:
		return attribute.Int64(key, v)
	case []int64:
		return attribute.Int64Slice(key, v)
	case float64:
		return attribute.Float64(key, v)
	case []float64:
		return attribute.Float64Slice(key, v)
	case string:
		return attribute.String(key, v)
	case []string:
		return attribute.StringSlice(key, v)
	default:
		return attribute.String(key, fmt.Sprintf("%v", v))
	}
}

func (s *Span) SetAttribute(key string, value interface{}) {
	s.Span.SetAttributes(generateAttribute(key, value))
}

func (s *Span) SetError(err error) {
	s.Span.RecordError(err)
}

func (s *Span) SetStatus(code sdk.Code, description string) {
	s.Span.SetStatus(codes.Code(code), description)
}

func (s *Span) IsNoop() bool {
	return !s.Span.IsRecording()
}

func (s *Span) AddEvent(name string, ts time.Time, attributes map[string]interface{}) {
	var otAttributes []attribute.KeyValue
	for k, v := range attributes {
		otAttributes = append(otAttributes, generateAttribute(k, v))
	}
	s.Span.AddEvent(name, trace.WithTimestamp(ts), trace.WithAttributes(otAttributes...))
}

func SpanFromContext(ctx context.Context) sdk.Span {
	return &Span{trace.SpanFromContext(ctx)}
}

type getTracerProvider func() trace.TracerProvider

func startSpan(provider getTracerProvider) sdk.StartSpan {
	return func(ctx context.Context, name string, opts *sdk.SpanOptions) (context.Context, sdk.Span, func()) {
		startOpts := []trace.SpanStartOption{}

		if opts != nil {
			startOpts = append(startOpts, trace.WithSpanKind(mapSpanKind(opts.Kind)))

			if opts.Timestamp.IsZero() {
				startOpts = append(startOpts, trace.WithTimestamp(time.Now()))
			} else {
				startOpts = append(startOpts, trace.WithTimestamp(opts.Timestamp))
			}
		}

		ctx, span := provider().
			Tracer(TracerDomain, trace.WithInstrumentationVersion(version.Version)).
			Start(ctx, name, startOpts...)
		return ctx, &Span{span}, func() { span.End() }
	}
}

var StartSpan = startSpan(otel.GetTracerProvider)
var NoopStartSpan = startSpan(trace.NewNoopTracerProvider)

func mapSpanKind(kind sdk.SpanKind) trace.SpanKind {
	switch kind {
	case sdk.SpanKindClient:
		return trace.SpanKindClient
	case sdk.SpanKindServer:
		return trace.SpanKindServer
	case sdk.SpanKindProducer:
		return trace.SpanKindProducer
	case sdk.SpanKindConsumer:
		return trace.SpanKindConsumer
	default:
		return trace.SpanKindUnspecified
	}
}
