# goagent for database/sql

A SQL wrapper attaching goagent instrumentation

## Usage

```go
import (
    "database/sql"
    "github.com/hypertrace/goagent/instrumentation/hypertrace/database/hypersql"
)

// Register our hypersql wrapper for the provided MySQL driver.
driverName, err = hypersql.Register("mysql")
if err != nil {
    log.Fatalf("unable to register goagent driver: %v\n", err)
}

// Connect to a MySQL database using the hypersql driver wrapper.
db, err = sql.Open(driverName, "user:password@/dbname")

```

You can also wrap your own driver with goagent instrumentation as follows:

```go

import (
    "github.com/go-sql-driver/mysql"
    "github.com/hypertrace/goagent/instrumentation/hypertrace/database/hypersql"
)

// Explicitly wrap the MySQL driver with hypersql
driver := hypersql.Wrap(&mysql.MySQLDriver{})

// Register our hypersql wrapper as a database driver
sql.Register("ht-mysql", driver)

// Connect to a MySQL database using the hypersql driver wrapper
db, err = sql.Open("ht-mysql", "user:password@/dbname")
```
