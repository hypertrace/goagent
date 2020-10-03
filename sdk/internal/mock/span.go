package mock

import (
	"context"

	"github.com/traceableai/goagent/sdk"
)

var _ sdk.Span = &Span{}

type Span struct {
	Attributes map[string]interface{}
	Noop       bool
}

func (s *Span) SetAttribute(key string, value interface{}) {
	if s.Attributes == nil {
		s.Attributes = make(map[string]interface{})
	}
	s.Attributes[key] = value
}

func (s *Span) IsNoop() bool {
	return s.Noop
}

type spanKey string

func SpanFromContext(ctx context.Context) sdk.Span {
	return ctx.Value(spanKey("span")).(*Span)
}

func ContextWithSpan(ctx context.Context, s sdk.Span) context.Context {
	return context.WithValue(ctx, spanKey("span"), s)
}
