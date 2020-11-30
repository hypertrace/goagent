package opencensus

import (
	"context"

	"github.com/hypertrace/goagent/sdk"
	"go.opencensus.io/trace"
)

var _ sdk.Span = &Span{nil}

type Span struct {
	*trace.Span
}

func (s *Span) IsNoop() bool {
	return !s.IsRecordingEvents()
}

func (s *Span) SetAttribute(key string, value interface{}) {
	switch v := value.(type) {
	case bool:
		s.Span.AddAttributes(trace.BoolAttribute(key, v))
	case int64:
		s.Span.AddAttributes(trace.Int64Attribute(key, v))
	case float64:
		s.Span.AddAttributes(trace.Float64Attribute(key, v))
	case string:
		s.Span.AddAttributes(trace.StringAttribute(key, v))
	}
}

func (s *Span) SetError(ctx context.Context, err error) {
	s.Span.AddAttributes(trace.StringAttribute("error", err.Error()))
}

func SpanFromContext(ctx context.Context) sdk.Span {
	return &Span{trace.FromContext(ctx)}
}

func StartSpan(ctx context.Context, name string, options *sdk.SpanOptions) (context.Context, sdk.Span, func()) {
	startOpts := []trace.StartOption{
		trace.WithSpanKind(toOCSpanKind(options.Kind)),
	}

	ctx, span := trace.StartSpan(ctx, name, startOpts...)
	return ctx, &Span{span}, span.End
}

func toOCSpanKind(kind sdk.SpanKind) int {
	switch kind {
	case sdk.Client:
		return trace.SpanKindClient
	case sdk.Server:
		return trace.SpanKindServer
	default:
		return trace.SpanKindUnspecified
	}
}
