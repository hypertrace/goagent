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
		return readbackSpanFromContext(ctx, v.(*map[string]interface{}))
	}

	return &Span{trace.SpanFromContext(ctx)}
}

func ReadbackSpanFromContext(ctx context.Context) sdk.ReadbackSpan {
	if v := ctx.Value(readAttrsKey); v != nil {
		return readbackSpanFromContext(ctx, v.(*map[string]interface{}))
	}

	return readbackSpanFromContext(ctx, &map[string]interface{}{})
}

func readbackSpanFromContext(ctx context.Context, attrs *map[string]interface{}) *ReadbackSpan {
	return &ReadbackSpan{&Span{trace.SpanFromContext(ctx)}, attrs, sync.Mutex{}}
}
