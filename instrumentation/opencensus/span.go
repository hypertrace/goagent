package opencensus

import (
	"context"
	"fmt"

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

func generateAttribute(key string, value interface{}) trace.Attribute {
	switch v := value.(type) {
	case bool:
		return trace.BoolAttribute(key, v)
	case int64:
		return trace.Int64Attribute(key, v)
	case float64:
		return trace.Float64Attribute(key, v)
	case string:
		return trace.StringAttribute(key, v)
	default:
		return trace.StringAttribute(key, fmt.Sprintf("%v", v))
	}
}

func (s *Span) SetAttribute(key string, value interface{}) {
	s.Span.AddAttributes(generateAttribute(key, value))
}

func (s *Span) SetError(err error) {
	s.Span.AddAttributes(trace.StringAttribute("error", err.Error()))
}

func SpanFromContext(ctx context.Context) sdk.Span {
	return &Span{trace.FromContext(ctx)}
}

func StartSpan(ctx context.Context, name string, options *sdk.SpanOptions) (context.Context, sdk.Span, func()) {
	startOpts := []trace.StartOption{
		trace.WithSpanKind(mapSpanKind(options.Kind)),
	}

	ctx, span := trace.StartSpan(ctx, name, startOpts...)
	return ctx, &Span{span}, span.End
}

func mapSpanKind(kind sdk.SpanKind) int {
	switch kind {
	case sdk.Client:
		return trace.SpanKindClient
	case sdk.Server:
		return trace.SpanKindServer
	default:
		return trace.SpanKindUnspecified
	}
}
