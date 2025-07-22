package hypersql // import "github.com/hypertrace/goagent/instrumentation/hypertrace/database/hypersql"

import (
	"database/sql/driver"
	otelsql "github.com/hypertrace/goagent/instrumentation/opentelemetry/database/hypersql"
)

// Wrap takes a SQL driver and wraps it with Hypertrace instrumentation.
func Wrap(d driver.Driver, opts ...Option) driver.Driver {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return otelsql.Wrap(d, o.toSDKOptions())
}

// Register initializes and registers the hypersql wrapped database driver
// identified by its driverName. On success it
// returns the generated driverName to use when calling sql.Open.
func Register(driverName string, opts ...Option) (string, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return otelsql.Register(driverName, o.toSDKOptions())

}
