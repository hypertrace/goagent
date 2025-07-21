package hypersql // import "github.com/hypertrace/goagent/instrumentation/opentelemetry/database/hypersql"

import (
	"database/sql/driver"
	"github.com/hypertrace/goagent/sdk/filter"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkSQL "github.com/hypertrace/goagent/sdk/instrumentation/database/sql"
)

// Wrap takes a SQL driver and wraps it with Hypertrace instrumentation.
func Wrap(d driver.Driver, filter filter.Filter) driver.Driver {
	return sdkSQL.Wrap(d, opentelemetry.StartSpan, filter)
}

// Register initializes and registers the hypersql wrapped database driver
// identified by its driverName. On success it
// returns the generated driverName to use when calling hypersql.Open.
func Register(driverName string, filter filter.Filter) (string, error) {
	return sdkSQL.Register(driverName, opentelemetry.StartSpan, filter)
}
