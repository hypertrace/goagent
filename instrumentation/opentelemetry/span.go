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

func StartSpan(ctx context.Context, name string) (context.Context, sdk.Span, func()) {
	ctx, span := global.Tracer("org.hypertrace.goagent").Start(ctx, name)
	return ctx, &Span{span}, func() { span.End() }
}
