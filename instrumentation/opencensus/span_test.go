package opencensus

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.opencensus.io/trace"

	"github.com/hypertrace/goagent/sdk"
	"github.com/stretchr/testify/assert"
)

func TestIsNoop(t *testing.T) {
	_, unsampledSpan := trace.StartSpan(context.Background(), "test", trace.WithSampler(trace.NeverSample()))
	span := &Span{unsampledSpan}
	assert.True(t, span.IsNoop())

	_, sampledSpan := trace.StartSpan(context.Background(), "test", trace.WithSampler(trace.AlwaysSample()))
	span = &Span{sampledSpan}
	assert.False(t, span.IsNoop())
}

func TestMapSpanKind(t *testing.T) {
	assert.Equal(t, mapSpanKind(sdk.SpanKindClient), trace.SpanKindClient)
	assert.Equal(t, mapSpanKind(sdk.SpanKindServer), trace.SpanKindServer)
}

func TestGenerateAttributeSuccess(t *testing.T) {
	const attrKey = "test_key"
	tCases := []struct {
		value        interface{}
		expectedAttr interface{}
	}{
		{value: true, expectedAttr: trace.BoolAttribute(attrKey, true)},
		{value: int64(1), expectedAttr: trace.Int64Attribute(attrKey, 1)},
		{value: float64(1.2), expectedAttr: trace.Float64Attribute(attrKey, 1.2)},
		{value: "abc", expectedAttr: trace.StringAttribute(attrKey, "abc")},
		{value: errors.New("xyz"), expectedAttr: trace.StringAttribute(attrKey, "xyz")},
	}

	for _, tCase := range tCases {
		assert.Equal(t, tCase.expectedAttr, generateAttribute(attrKey, tCase.value))
	}
}

func TestAddEvents(t *testing.T) {
	_, s, _ := StartSpan(context.Background(), "test_span", &sdk.SpanOptions{})
	m := make(map[string]interface{})
	s.AddEvent("test_event_1", time.Now(), m)
	m["k1"] = "v1"
	s.AddEvent("test_event_2", time.Now(), m)
	m["k2"] = 23
	s.AddEvent("test_event_3", time.Now(), m)
	m["k3"] = true
	s.AddEvent("test_event_4", time.Now(), m)
}
