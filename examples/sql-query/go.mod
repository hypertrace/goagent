module github.com/hypertrace/goagent/examples/sql-query

go 1.15

replace github.com/hypertrace/goagent => ../..

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/hypertrace/goagent v0.0.0-00010101000000-000000000000
)
