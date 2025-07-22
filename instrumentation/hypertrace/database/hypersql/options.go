package hypersql // import "github.com/hypertrace/goagent/instrumentation/hypertrace/database/hypersql"

import (
	"github.com/hypertrace/goagent/sdk/filter"
	sdkSQL "github.com/hypertrace/goagent/sdk/instrumentation/database/sql"
)

type options struct {
	Filter filter.Filter
}

func (o *options) toSDKOptions() *sdkSQL.Options {
	opts := (sdkSQL.Options)(*o)
	return &opts
}

type Option func(o *options)

// WithFilter adds a filter to the GRPC option.
func WithFilter(f filter.Filter) Option {
	return func(o *options) {
		o.Filter = f
	}
}
