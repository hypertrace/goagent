package opentelemetry_test

import (
	"context"
	"testing"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
)

func TestIsNoop(t *testing.T) {
	span := &opentelemetry.Span{trace.NoopSpan{}}
	assert.True(t, span.IsNoop())

	internal.InitTracer()
	_, delegateSpan := global.TraceProvider().Tracer("ai.traceable").Start(context.Background(), "test_span")
	span = &opentelemetry.Span{delegateSpan}
	assert.False(t, span.IsNoop())
}
