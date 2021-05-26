package opentelemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadbackSpanFromContextForEmptyContext(t *testing.T) {
	_, s := ReadbackSpanFromContext(context.Background())
	assert.Empty(t, s.GetAttributes())
}

func TestReadbackSpanFromContext(t *testing.T) {
	ctx, s := ReadbackSpanFromContext(context.Background())
	assert.True(t, s.IsNoop())

	s.SetAttribute("http.method", "GET")
	assert.Equal(t, "GET", s.GetAttributes()["http.method"].(string))
	ctxAttrs := ctx.Value(readAttrsKey).(*map[string]interface{})
	assert.Equal(t, "GET", (*ctxAttrs)["http.method"].(string))

	es := SpanFromContext(ctx)
	assert.Equal(t, s, es)
}
