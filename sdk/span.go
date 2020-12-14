package sdk

import "context"

// Span is an interface that accept attributed and can be
// distinguished as noop
type Span interface {
	// SetAttribute sets an attribute for the span.
	SetAttribute(key string, value interface{})

	// SetError sets an error for the span.
	SetError(err error)

	// IsNoop tells whether the span is noop or not, useful for avoiding
	// expensive recording.
	IsNoop() bool
}

// SpanFromContext retrieves the existing span from a context
type SpanFromContext func(ctx context.Context) Span

// SpanKind represents the span kind, either Client, Server,
// Producer, Consumer or undertermine
type SpanKind string

const (
	Undetermined SpanKind = ""
	Client       SpanKind = "CLIENT"
	Server       SpanKind = "SERVER"
	Producer     SpanKind = "PRODUCER"
	Consumer     SpanKind = "CONSUMER"
)

// SpanOptions describes the options for starting a span
type SpanOptions struct {
	Kind SpanKind
}

// StartSpan creates a span and injects into a context, returning a span ender function
type StartSpan func(ctx context.Context, name string, opts *SpanOptions) (context.Context, Span, func())
