package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopFilter(t *testing.T) {
	f := NoopFilter{}
	res := f.Evaluate(nil)
	assert.False(t, res.Block)
}
