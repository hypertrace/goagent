package http

import (
	"net/http"
	"testing"

	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestSetScalarAttributeSuccess(t *testing.T) {
	h := http.Header{}
	h.Set("key_1", "value_1")
	span := mock.NewSpan()
	setAttributesFromHeaders("request", h, span)
	assert.Equal(t, "value_1", span.ReadAttribute("http.request.header.key_1").(string))

	_ = span.ReadAttribute("container_id") // needed in containarized envs
	assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
}

func TestSetMultivalueAttributeSuccess(t *testing.T) {
	h := http.Header{}
	h.Add("key_1", "value_1")
	h.Add("key_1", "value_2")

	span := mock.NewSpan()
	setAttributesFromHeaders("response", h, span)

	assert.Equal(t, "value_1", span.ReadAttribute("http.response.header.key_1[0]").(string))
	assert.Equal(t, "value_2", span.ReadAttribute("http.response.header.key_1[1]").(string))

	_ = span.ReadAttribute("container_id") // needed in containarized envs
	assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
}
