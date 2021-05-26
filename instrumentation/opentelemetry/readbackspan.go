package opentelemetry

import (
	"context"
	"sync"

	"github.com/hypertrace/goagent/sdk"
	"go.opentelemetry.io/otel"
)

var _ sdk.ReadbackSpan = (*ReadbackSpan)(nil)

type ReadbackSpan struct {
	sdk.Span
	readAttrs *map[string]interface{}
	m         sync.Mutex
}

func (rs *ReadbackSpan) GetAttributes() map[string]interface{} {
	rs.m.Lock()
	defer rs.m.Unlock()

	return *rs.readAttrs
}

func (rs *ReadbackSpan) SetAttribute(key string, value interface{}) {
	rs.m.Lock()
	defer rs.m.Unlock()

	(*rs.readAttrs)[key] = value
	rs.Span.SetAttribute(key, value)
}

func startReadbackSpan(provider getTracerProvider) sdk.StartReadbackSpan {
	ss := startSpan(provider)
	return func(ctx context.Context, name string, opts *sdk.SpanOptions) (context.Context, sdk.ReadbackSpan, func()) {
		ctx, s, ender := ss(ctx, name, opts)
		readAttrs := map[string]interface{}{}
		return context.WithValue(ctx, readAttrsKey, &readAttrs), &ReadbackSpan{s, &readAttrs, sync.Mutex{}}, ender
	}
}

var StartReadbackSpan = startReadbackSpan(otel.GetTracerProvider)
