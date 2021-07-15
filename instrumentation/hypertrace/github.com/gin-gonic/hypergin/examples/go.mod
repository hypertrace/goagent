module github.com/hypertrace/goagent/instrumentation/hypertrace/github.com/gin-gonic/gypergin/examples/main

go 1.16

replace github.com/hypertrace/goagent => ../../../../../../

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/hypertrace/goagent v0.2.1
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.18.0 // indirect
)
