// +build ignore

package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/database/hypersql"
)

// Run docker run mysql
func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("querier")

	flusher := hypertrace.Init(cfg)
	defer flusher()

	var (
		driver driver.Driver
		db     *sql.DB
	)

	// Explicitly wrap the MySQLDriver driver with hypersql.
	driver = hypersql.Wrap(&mysql.MySQLDriver{})

	// Register our hypersql wrapper as a database driver.
	sql.Register("ht-mysql", driver)

	db, err := sql.Open("ht-mysql", "root:root@tcp(localhost)/")
	if err != nil {
		log.Fatalf("failed to connect the DB: %v", err)
	}

	const dbPingRetries = 5
	for i := 0; i <= dbPingRetries; i++ {
		if err := db.Ping(); err != nil && i == dbPingRetries {
			log.Fatalf("failed to ping the DB: %v", err)
		}
		time.Sleep(time.Second)
	}

	rows, err := db.QueryContext(context.Background(), "SELECT 'Hi there :)'")
	if err != nil {
		log.Fatalf("failed to retrieve message: %v", err)
	}

	for rows.Next() {
		m := new(string)
		rows.Scan(m)
		fmt.Println(*m)
	}
}
