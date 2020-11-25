package opentelemetry_test

import (
	"context"
	"testing"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
)

func TestIsNoop(t *testing.T) {
	span := &opentelemetry.Span{trace.NoopSpan{}}
	assert.True(t, span.IsNoop())

	opentelemetry.Init(config.Load())
	_, delegateSpan := global.Tracer(opentelemetry.TracerDomain).Start(context.Background(), "test_span")
	span = &opentelemetry.Span{delegateSpan}
	assert.False(t, span.IsNoop())
}
