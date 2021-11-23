package opentelemetry

import (
	"context"
	"errors"
	"testing"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/sdk"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func newNoopSpan() trace.Span {
	_, noopSpan := trace.NewNoopTracerProvider().Tracer("noop").Start(context.Background(), "test_name")
	return noopSpan
}

func TestIsNoop(t *testing.T) {
	span := &Span{newNoopSpan()}
	assert.True(t, span.IsNoop())

	shutdown := Init(config.Load())
	defer shutdown()

	_, delegateSpan := otel.Tracer(TracerDomain).Start(context.Background(), "test_span")
	span = &Span{delegateSpan}
	assert.False(t, span.IsNoop())
}

func TestMapSpanKind(t *testing.T) {
	assert.Equal(t, trace.SpanKindClient, mapSpanKind(sdk.SpanKindClient))
	assert.Equal(t, trace.SpanKindServer, mapSpanKind(sdk.SpanKindServer))
	assert.Equal(t, trace.SpanKindProducer, mapSpanKind(sdk.SpanKindProducer))
	assert.Equal(t, trace.SpanKindConsumer, mapSpanKind(sdk.SpanKindConsumer))
}

func TestSetAttributeSuccess(t *testing.T) {
	_, s, _ := StartSpan(context.Background(), "test_span", &sdk.SpanOptions{})
	s.SetAttribute("test_key_1", true)
	s.SetAttribute("test_key_2", int64(1))
	s.SetAttribute("test_key_3", float64(1.2))
	s.SetAttribute("test_key_4", "abc")
	s.SetAttribute("test_key_4", errors.New("xyz"))
}

func TestGenerateAttribute(t *testing.T) {
	assert.Equal(t, attribute.BOOL, generateAttribute("key", true).Value.Type())
	assert.Equal(t, attribute.BOOLSLICE, generateAttribute("key", []bool{true}).Value.Type())
	assert.Equal(t, attribute.INT64, generateAttribute("key", 1).Value.Type())
	assert.Equal(t, attribute.INT64SLICE, generateAttribute("key", []int{1}).Value.Type())
	assert.Equal(t, attribute.INT64, generateAttribute("key", int64(1)).Value.Type())
	assert.Equal(t, attribute.INT64SLICE, generateAttribute("key", []int64{1}).Value.Type())
	assert.Equal(t, attribute.FLOAT64, generateAttribute("key", 1.23).Value.Type())
	assert.Equal(t, attribute.FLOAT64SLICE, generateAttribute("key", []float64{1.23}).Value.Type())
	assert.Equal(t, attribute.STRING, generateAttribute("key", "val").Value.Type())
	assert.Equal(t, attribute.STRINGSLICE, generateAttribute("key", []string{"val"}).Value.Type())

	attr := generateAttribute("key", errors.New("x"))
	assert.Equal(t, attribute.STRING, attr.Value.Type())
	assert.Equal(t, "x", attr.Value.AsString())
}
