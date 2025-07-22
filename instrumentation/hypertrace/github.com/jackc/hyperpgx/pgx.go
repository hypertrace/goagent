package hyperpgx // import "github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/jackc/hyperpgx"

import (
	"context"

	otelpgx "github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/jackc/hyperpgx"
)

func Connect(ctx context.Context, connString string, opts ...Option) (otelpgx.PGXConn, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return otelpgx.Connect(ctx, connString, o.toSDKOptions())
}
