package hypergin // import "github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/gin-gonic/hypergin"

import (
	"github.com/gin-gonic/gin"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/gin-gonic/hypergin"
)

func Middleware(opts ...Option) gin.HandlerFunc {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return hypergin.Middleware(o.toSDKOptions())
}
