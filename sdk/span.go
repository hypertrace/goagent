package sdk // import "github.com/hypertrace/goagent/sdk"

import (
	"context"
	"time"
)

type AttributeList interface {
	GetValue(key string) interface{}

	// Iterate loops through the attributes list and applies the yield function on each attribute.
	// If the yield function returns false, we exit the loop.
	Iterate(yield func(key string, value interface{}) bool)

	Len() int
}

// Span is an interface that accepts attributes and can be
// distinguished as noop
type Span interface {
	GetAttributes() AttributeList

	// SetAttribute sets an attribute for the span.
	SetAttribute(key string, value interface{})

	// SetError sets an error for the span.
	SetError(err error)

	// SetStatus sets the status of the Span in the form of a code and a
	// description.
	SetStatus(code Code, description string)

	// IsNoop tells whether the span is noop or not, useful for avoiding
	// expensive recording.
	IsNoop() bool

	// AddEvent adds an event to the Span with the provided name, timestamp and attributes.
	AddEvent(name string, ts time.Time, attributes map[string]interface{})
}

// SpanFromContext retrieves the existing span from a context
type SpanFromContext func(ctx context.Context) Span

// SpanKind represents the span kind, either Client, Server,
// Producer, Consumer or undertermine
type SpanKind string

const (
	SpanKindUndetermined SpanKind = ""
	SpanKindClient       SpanKind = "CLIENT"
	SpanKindServer       SpanKind = "SERVER"
	SpanKindProducer     SpanKind = "PRODUCER"
	SpanKindConsumer     SpanKind = "CONSUMER"
)

// SpanOptions describes the options for starting a span
type SpanOptions struct {
	Kind      SpanKind
	Timestamp time.Time
}

// StartSpan creates a span and injects into a context, returning a span ender function
type StartSpan func(ctx context.Context, name string, opts *SpanOptions) (context.Context, Span, func())
