module github.com/hypertrace/goagent/examples/postgres-query

go 1.15

replace github.com/hypertrace/goagent => ../..

replace github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/jackc/hyperpgx => ../../instrumentation/opentelemetry/github.com/jackc/hyperpgx

require github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/jackc/hyperpgx v0.0.0-00010101000000-000000000000
