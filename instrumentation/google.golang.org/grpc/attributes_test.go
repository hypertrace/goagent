package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traceableai/goagent/instrumentation/internal"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/grpc/metadata"
)

func TestSetScalarAttributeSuccess(t *testing.T) {
	_, flusher := internal.InitTracer()

	md := metadata.Pairs("key_1", "value_1")
	_, span := global.TraceProvider().Tracer("ai.traceable").Start(context.Background(), "")
	setAttributes(md, span)
	span.End()

	readbackSpan := flusher()[0]
	attrs := internal.LookupAttributes(readbackSpan.Attributes)
	assert.Equal(t, "value_1", attrs.Get("grpc.request.metadata.key_1").AsString())
}

func TestSetMultivalueAttributeSuccess(t *testing.T) {
	_, flusher := internal.InitTracer()

	md := metadata.Pairs("key_1", "value_1", "key_1", "value_2")
	_, span := global.TraceProvider().Tracer("ai.traceable").Start(context.Background(), "")
	setAttributes(md, span)
	span.End()

	readbackSpan := flusher()[0]
	attrs := internal.LookupAttributes(readbackSpan.Attributes)
	assert.Equal(t, "value_1", attrs.Get("grpc.request.metadata.key_1[0]").AsString())
	assert.Equal(t, "value_2", attrs.Get("grpc.request.metadata.key_1[1]").AsString())
}
