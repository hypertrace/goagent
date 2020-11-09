package hypersql

import (
	"database/sql/driver"

	"contrib.go.opencensus.io/integrations/ocsql"
)

// Wrap takes a SQL driver and wraps it with Hypertrace instrumentation.
func Wrap(d driver.Driver) driver.Driver {
	return ocsql.Wrap(d,
		ocsql.WithQuery(true),
		ocsql.WithDisableErrSkip(true),
	)
}

// Register initializes and registers the hypersql wrapped database driver
// identified by its driverName. On success it
// returns the generated driverName to use when calling sql.Open.
func Register(driverName string) (string, error) {
	return ocsql.Register(
		driverName,
		ocsql.WithQuery(true),
		ocsql.WithDisableErrSkip(true),
	)
}
