module github.com/hypertrace/goagent/examples/mux-server

go 1.15

replace github.com/hypertrace/goagent => ../..

require (
	github.com/gorilla/mux v1.8.0
	github.com/hypertrace/goagent v0.0.0-00010101000000-000000000000
)
