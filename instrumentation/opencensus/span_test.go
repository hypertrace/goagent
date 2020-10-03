package opencensus_test

import (
	"context"
	"testing"

	"go.opencensus.io/trace"

	"github.com/stretchr/testify/assert"
	"github.com/traceableai/goagent/instrumentation/opencensus"
)

func TestIsNoop(t *testing.T) {
	_, unsampledSpan := trace.StartSpan(context.Background(), "test", trace.WithSampler(trace.NeverSample()))
	span := &opencensus.Span{unsampledSpan}
	assert.True(t, span.IsNoop())

	_, sampledSpan := trace.StartSpan(context.Background(), "test", trace.WithSampler(trace.AlwaysSample()))
	span = &opencensus.Span{sampledSpan}
	assert.False(t, span.IsNoop())
}
