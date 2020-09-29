package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traceableai/goagent/sdk/internal/mock"
)

func TestSetScalarAttributeSuccess(t *testing.T) {
	h := http.Header{}
	h.Set("key_1", "value_1")
	span := &mock.Span{}
	setAttributesFromHeaders("request", h, span)
	assert.Equal(t, "value_1", span.Attributes["http.request.header.Key_1"].(string))
}

func TestSetMultivalueAttributeSuccess(t *testing.T) {
	h := http.Header{}
	h.Add("key_1", "value_1")
	h.Add("key_1", "value_2")

	span := &mock.Span{}
	setAttributesFromHeaders("response", h, span)

	assert.Equal(t, "value_1", span.Attributes["http.response.header.Key_1[0]"].(string))
	assert.Equal(t, "value_2", span.Attributes["http.response.header.Key_1[1]"].(string))
}
