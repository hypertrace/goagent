package hypersql // import "github.com/hypertrace/goagent/instrumentation/hypertrace/database/hypersql"

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
