package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupSuccess(t *testing.T) {
	h := http.Header{}
	h.Set("abc", "123")

	a := headerMapAccessor{h}
	assert.Equal(t, a.Lookup("abc"), []string{"123"})
	assert.Equal(t, a.Lookup("aBC"), []string{"123"})
	assert.Empty(t, a.Lookup("xyz"))
}
