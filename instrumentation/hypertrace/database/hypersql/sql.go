package hypersql

import (
	otelsql "github.com/hypertrace/goagent/instrumentation/opentelemetry/database/hypersql"
)

// Wrap takes a SQL driver and wraps it with Hypertrace instrumentation.
var Wrap = otelsql.Wrap

// Register initializes and registers the hypersql wrapped database driver
// identified by its driverName. On success it
// returns the generated driverName to use when calling sql.Open.
var Register = otelsql.Register
