package opentelemetry

import (
	"context"

	"github.com/hypertrace/goagent/sdk"
	"go.opentelemetry.io/otel/api/trace"
)

var _ sdk.Span = &Span{trace.NoopSpan{}}

type Span struct {
	trace.Span
}

func (s *Span) IsNoop() bool {
	_, ok := (s.Span).(trace.NoopSpan)
	return ok
}

func SpanFromContext(ctx context.Context) sdk.Span {
	return &Span{trace.SpanFromContext(ctx)}
}
