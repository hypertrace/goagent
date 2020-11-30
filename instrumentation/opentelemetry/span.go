package opentelemetry

import (
	"context"

	"github.com/hypertrace/goagent/sdk"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
)

var _ sdk.Span = &Span{trace.NoopSpan{}}

type Span struct {
	trace.Span
}

func (s *Span) SetError(ctx context.Context, err error) {
	s.Span.RecordError(ctx, err)
}

func (s *Span) IsNoop() bool {
	_, ok := (s.Span).(trace.NoopSpan)
	return ok
}

func SpanFromContext(ctx context.Context) sdk.Span {
	return &Span{trace.SpanFromContext(ctx)}
}

func StartSpan(ctx context.Context, name string, options *sdk.SpanOptions) (context.Context, sdk.Span, func()) {
	startOpts := []trace.StartOption{
		trace.WithSpanKind(toOTelSpanKind(options.Kind)),
	}

	ctx, span := global.Tracer(TracerDomain).Start(ctx, name, startOpts...)
	return ctx, &Span{span}, func() { span.End() }
}

func toOTelSpanKind(kind sdk.SpanKind) trace.SpanKind {
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
