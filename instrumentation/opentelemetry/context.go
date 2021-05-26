package opentelemetry

import (
	"context"
	"sync"

	"github.com/hypertrace/goagent/sdk"
	"go.opentelemetry.io/otel/trace"
)

var readAttrsKey = struct{}{}

func SpanFromContext(ctx context.Context) sdk.Span {
	if v := ctx.Value(readAttrsKey); v != nil {
		s := trace.SpanFromContext(ctx)
		return readbackSpanFromContext(s, v.(*map[string]interface{}))
	}

	return &Span{trace.SpanFromContext(ctx)}
}

// ReadbackSpanFromContext extracts the span from the context if it exists or
// put a noop span and return a new context including it.
func ReadbackSpanFromContext(ctx context.Context) (context.Context, sdk.ReadbackSpan) {
	if v := ctx.Value(readAttrsKey); v != nil {
		s := trace.SpanFromContext(ctx)
		return ctx, readbackSpanFromContext(s, v.(*map[string]interface{}))
	}

	s := trace.SpanFromContext(ctx)
	readAttrs := map[string]interface{}{}
	// even if the span from the context is noop, we need to store the attributes
	// and make sure the span is passed down in the context because it could be read
	// later.
	return context.WithValue(trace.ContextWithSpan(ctx, s), readAttrsKey, &readAttrs),
		readbackSpanFromContext(s, &readAttrs)
}

func readbackSpanFromContext(s trace.Span, attrs *map[string]interface{}) *ReadbackSpan {
	return &ReadbackSpan{&Span{s}, attrs, sync.Mutex{}}
}
