package opentelemetry

import (
	"context"

	"github.com/hypertrace/goagent/sdk"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

var _ sdk.Span = (*Span)(nil)

type Span struct {
	trace.Span
}

func (s *Span) SetAttribute(key string, value interface{}) {
	s.Span.SetAttributes(label.Any(key, value))
}

func (s *Span) SetError(err error) {
	s.Span.RecordError(err)
}

func (s *Span) IsNoop() bool {
	return !s.Span.IsRecording()
}

func SpanFromContext(ctx context.Context) sdk.Span {
	return &Span{trace.SpanFromContext(ctx)}
}

func StartSpan(ctx context.Context, name string, options *sdk.SpanOptions) (context.Context, sdk.Span, func()) {
	startOpts := []trace.SpanOption{
		trace.WithSpanKind(mapSpanKind(options.Kind)),
	}

	ctx, span := otel.Tracer(TracerDomain).Start(ctx, name, startOpts...)
	return ctx, &Span{span}, func() { span.End() }
}

func mapSpanKind(kind sdk.SpanKind) trace.SpanKind {
	switch kind {
	case sdk.Client:
		return trace.SpanKindClient
	case sdk.Server:
		return trace.SpanKindServer
	case sdk.Producer:
		return trace.SpanKindProducer
	case sdk.Consumer:
		return trace.SpanKindConsumer
	default:
		return trace.SpanKindUnspecified
	}
}
