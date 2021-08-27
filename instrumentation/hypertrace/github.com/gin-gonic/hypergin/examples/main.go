package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/gin-gonic/hypergin"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// hypergin.DisableConsoleColor()
	r := gin.Default()

	r.Use(hypergin.Middleware(&sdkhttp.Options{}))

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    http.StatusOK,
			"message": "pong",
		})
	})

	return r
}

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("gin-example-server")

	flusher := hypertrace.Init(cfg)
	defer flusher()

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
