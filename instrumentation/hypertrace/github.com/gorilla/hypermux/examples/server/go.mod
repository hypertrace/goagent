module github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/gorilla/hypermux/examples/server

go 1.16

replace github.com/hypertrace/goagent => ../../../../../../../

//replace github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/gorilla/hypermux => ../../../../../../opentelemetry/github.com/gorilla/hypermux

require (
	github.com/gorilla/mux v1.8.0
	github.com/hypertrace/goagent v0.0.0-00010101000000-000000000000
)
