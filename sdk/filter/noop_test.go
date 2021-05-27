package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopFilter(t *testing.T) {
	f := NoopFilter{}
	assert.False(t, f.EvaluateURLAndHeaders(nil, "", nil))
	assert.False(t, f.EvaluateBody(nil, nil))
	assert.False(t, f.Evaluate(nil))
}
