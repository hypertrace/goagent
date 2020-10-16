package sdk

import "context"

// Span is an interface that accept attributed and can be
// distinguished as noop
type Span interface {
	// SetAttribute sets an attribute for the span.
	SetAttribute(key string, value interface{})

	// IsNoop tells whether the span is noop or not, useful for avoiding
	// expensive recording.
	IsNoop() bool
}

// SpanFromContext retrieves the existing span from a context
type SpanFromContext func(ctx context.Context) Span
