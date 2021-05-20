package http

import (
	"testing"

	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestBodyTruncationSuccess(t *testing.T) {
	s := mock.NewSpan()
	setTruncatedBodyAttribute("request", []byte("text"), 2, s)
	assert.Equal(t, "te", s.ReadAttribute("http.request.body"))
	assert.True(t, (s.ReadAttribute("http.request.body.truncated")).(bool))
	assert.Zero(t, s.RemainingAttributes())
}

func TestBodyTruncationIsSkipped(t *testing.T) {
	s := mock.NewSpan()
	setTruncatedBodyAttribute("request", []byte("text"), 7, s)
	assert.Equal(t, "text", s.ReadAttribute("http.request.body"))
	assert.Zero(t, s.RemainingAttributes())
}
