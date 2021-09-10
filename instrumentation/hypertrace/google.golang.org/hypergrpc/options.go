package hypergrpc // import "github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc"

import (
	"github.com/hypertrace/goagent/sdk/filter"
	"github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"
)

type options struct {
	Filter filter.Filter
}

func (o *options) toSDKOptions() *grpc.Options {
	opts := (grpc.Options)(*o)
	return &opts
}

type Option func(o *options)

// WithFilter adds a filter to the GRPC option.
func WithFilter(f filter.Filter) Option {
	return func(o *options) {
		o.Filter = f
	}
}
