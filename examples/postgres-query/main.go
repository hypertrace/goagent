package main

// Copied from https://github.com/jackc/pgx#example-usage

import (
	"context"
	"fmt"
	"os"

	// gosec complains about this pkg not following golang repo standards
	// "could not import github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/jackc/hyperpgx (invalid package name: "")"
	// It is caused the pkg having its own go.mod
	"github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/jackc/hyperpgx" // #nosec
)

func main() {
	conn, err := hyperpgx.Connect(context.Background(), "root:root@localhost")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	m := new(string)
	err = conn.QueryRow(context.Background(), "SELECT 'Hi there :)'").Scan(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(*m)
}
