package sdk

import "context"

// Span is an interface that accept attributed and can be
// distinguished as noop
type Span interface {
	SetAttribute(key string, value interface{})
	IsNoop() bool
}

// SpanFromContext retrieves the existing span from a context
type SpanFromContext func(ctx context.Context) Span
