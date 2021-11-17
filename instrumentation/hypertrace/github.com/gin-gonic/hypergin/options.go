package hypergin // import "github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/gin-gonic/hypergin"

import (
	"github.com/hypertrace/goagent/sdk/filter"
	"github.com/hypertrace/goagent/sdk/instrumentation/net/http"
)

type options struct {
	Filter filter.Filter
}

func (o *options) toSDKOptions() *http.Options {
	opts := (http.Options)(*o)
	return &opts
}

type Option func(o *options)

func WithFilter(f filter.Filter) Option {
	return func(o *options) {
		o.Filter = f
	}
}
