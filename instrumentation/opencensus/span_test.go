package opencensus

import (
	"context"
	"testing"

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
	assert.Equal(t, mapSpanKind(sdk.Client), trace.SpanKindClient)
	assert.Equal(t, mapSpanKind(sdk.Server), trace.SpanKindServer)
}
