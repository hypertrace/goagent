package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopFilter(t *testing.T) {
	f := NoopFilter{}
	res, _ := f.EvaluateURLAndHeaders(nil, "", nil)
	assert.False(t, res)
	res, _ = f.EvaluateBody(nil, nil, nil)
	assert.False(t, res)
	res, _ = f.Evaluate(nil, "", nil, nil)
	assert.False(t, res)
}
