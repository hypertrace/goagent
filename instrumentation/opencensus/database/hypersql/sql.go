package hypersql // import "github.com/hypertrace/goagent/instrumentation/opencensus/database/hypersql"

import (
	"database/sql/driver"

	"github.com/hypertrace/goagent/instrumentation/opencensus"
	sdkSQL "github.com/hypertrace/goagent/sdk/instrumentation/database/sql"
)

// Wrap takes a SQL driver and wraps it with Hypertrace instrumentation.
func Wrap(d driver.Driver) driver.Driver {
	return sdkSQL.Wrap(d, opencensus.StartSpan)
}

// Register initializes and registers the hypersql wrapped database driver
// identified by its driverName. On success it
// returns the generated driverName to use when calling sql.Open.
func Register(driverName string) (string, error) {
	return sdkSQL.Register(driverName, opencensus.StartSpan)
}
