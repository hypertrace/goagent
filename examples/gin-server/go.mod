module github.com/hypertrace/goagent/examples/gin-server

go 1.15

replace github.com/hypertrace/goagent => ../..

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/hypertrace/goagent v0.2.1
)