package hyperpgx // import "github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/jackc/hyperpgx"

import (
	otelpgx "github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/jackc/hyperpgx"
	"github.com/hypertrace/goagent/sdk/filter"
)

type options struct {
	Filter filter.Filter
}

func (o *options) toSDKOptions() *otelpgx.Options {
	opts := (otelpgx.Options)(*o)
	return &opts
}

type Option func(o *options)

// WithFilter adds a filter to the GRPC option.
func WithFilter(f filter.Filter) Option {
	return func(o *options) {
		o.Filter = f
	}
}
