package opentelemetry

import (
	"context"
	"testing"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/sdk"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func newNoopSpan() trace.Span {
	_, noopSpan := trace.NewNoopTracerProvider().Tracer("noop").Start(context.Background(), "test_name")
	return noopSpan
}

func TestIsNoop(t *testing.T) {
	span := &Span{newNoopSpan()}
	assert.True(t, span.IsNoop())

	Init(config.Load())
	_, delegateSpan := otel.Tracer(TracerDomain).Start(context.Background(), "test_span")
	span = &Span{delegateSpan}
	assert.False(t, span.IsNoop())
}

func TestMapSpanKind(t *testing.T) {
	assert.Equal(t, mapSpanKind(sdk.Client), trace.SpanKindClient)
	assert.Equal(t, mapSpanKind(sdk.Server), trace.SpanKindServer)
	assert.Equal(t, mapSpanKind(sdk.Producer), trace.SpanKindProducer)
	assert.Equal(t, mapSpanKind(sdk.Consumer), trace.SpanKindConsumer)
}
