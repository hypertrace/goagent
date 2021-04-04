package main

// Copied from https://github.com/jackc/pgx#example-usage

import (
	"context"
	"fmt"
	"os"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/jackc/hyperpgx"
	"github.com/jackc/pgx/v4"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	wConn := hyperpgx.WrapConnection(conn)

	var name string
	var weight int64
	err = wConn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)
}
