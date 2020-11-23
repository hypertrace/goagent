package hypersql

import (
	"database/sql/driver"

	"github.com/hypertrace/goagent/instrumentation/opencensus"
	sqlsdk "github.com/hypertrace/goagent/sdk/database/hypersql"
)

// Wrap takes a SQL driver and wraps it with Hypertrace instrumentation.
func Wrap(d driver.Driver) driver.Driver {
	return sqlsdk.Wrap(d, opencensus.StartSpan)
}

// Register initializes and registers the hypersql wrapped database driver
// identified by its driverName. On success it
// returns the generated driverName to use when calling sql.Open.
func Register(driverName string) (string, error) {
	return sqlsdk.Register(driverName, opencensus.StartSpan)
}
