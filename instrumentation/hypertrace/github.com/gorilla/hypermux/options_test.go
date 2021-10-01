package hypermux // import "github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/gorilla/hypermux"

import (
	"testing"

	"github.com/hypertrace/goagent/sdk/filter"
	"github.com/stretchr/testify/assert"
)

func TestOptionsToSDK(t *testing.T) {
	o := &options{
		Filter: filter.NoopFilter{},
	}
	assert.Equal(t, filter.NoopFilter{}, o.toSDKOptions().Filter)
}
