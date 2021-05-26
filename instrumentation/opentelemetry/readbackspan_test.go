package opentelemetry

import (
	"context"
	"errors"
	"testing"

	"github.com/hypertrace/goagent/sdk"
	"github.com/stretchr/testify/assert"
)

func TestSetAndGetAttributeSuccess(t *testing.T) {
	ctx, s, _ := StartReadbackSpan(context.Background(), "test_readbackspan", &sdk.SpanOptions{})
	s.SetAttribute("test_key_1", true)
	s.SetAttribute("test_key_2", int64(1))
	s.SetAttribute("test_key_3", float64(1.2))
	s.SetAttribute("test_key_4", "abc")
	err := errors.New("xyz")
	s.SetAttribute("test_key_4", err)
	attributes := s.GetAttributes()
	assert.Equal(t, attributes["test_key_1"].(bool), true)
	assert.Equal(t, attributes["test_key_2"].(int64), int64(1))
	assert.Equal(t, attributes["test_key_3"].(float64), 1.2)
	assert.Equal(t, attributes["test_key_4"].(error), err)

	ctxAttributes := ctx.Value(readAttrsKey).(*map[string]interface{})
	assert.Equal(t, (*ctxAttributes)["test_key_1"].(bool), true)
	assert.Equal(t, (*ctxAttributes)["test_key_2"].(int64), int64(1))
	assert.Equal(t, (*ctxAttributes)["test_key_3"].(float64), 1.2)
	assert.Equal(t, (*ctxAttributes)["test_key_4"].(error), err)
}
