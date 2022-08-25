package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopFilter(t *testing.T) {
	f := NoopFilter{}
	res := f.EvaluateURLAndHeaders(nil, "", nil)
	assert.False(t, res.Block)
	res = f.EvaluateBody(nil, nil, nil)
	assert.False(t, res.Block)
	res = f.Evaluate(nil, "", nil, nil)
	assert.False(t, res.Block)
}
