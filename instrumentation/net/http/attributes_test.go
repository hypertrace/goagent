package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traceableai/goagent/instrumentation/internal"
	"go.opentelemetry.io/otel/api/global"
)

func TestSetScalarAttributeSuccess(t *testing.T) {
	_, flusher := internal.InitTracer()

	h := http.Header{}
	h.Set("key_1", "value_1")
	_, span := global.TraceProvider().Tracer("ai.traceable").Start(context.Background(), "")
	setAttributesFromHeaders("request", h, span)
	span.End()

	readbackSpan := flusher()[0]
	attrs := internal.LookupAttributes(readbackSpan.Attributes)
	assert.Equal(t, "value_1", attrs.Get("http.request.header.Key_1").AsString())
}

func TestSetMultivalueAttributeSuccess(t *testing.T) {
	_, flusher := internal.InitTracer()

	h := http.Header{}
	h.Add("key_1", "value_1")
	h.Add("key_1", "value_2")

	_, span := global.TraceProvider().Tracer("ai.traceable").Start(context.Background(), "")
	setAttributesFromHeaders("response", h, span)
	span.End()

	readbackSpan := flusher()[0]
	attrs := internal.LookupAttributes(readbackSpan.Attributes)
	assert.Equal(t, "value_1", attrs.Get("http.response.header.Key_1[0]").AsString())
	assert.Equal(t, "value_2", attrs.Get("http.response.header.Key_1[1]").AsString())
}
