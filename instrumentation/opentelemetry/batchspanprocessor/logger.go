package batchspanprocessor // import "github.com/hypertrace/goagent/instrumentation/opentelemetry/batchspanprocessor"

// Adapted from go.opentelemetry.io/otel/internal/global#internal_logging.go
import (
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

// The logger uses stdr which is backed by the standard `log.Logger`
// interface. This logger will only show messages at the Error Level.
var logger logr.Logger = stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))

// Info prints messages about the general state of the API or SDK.
// This should usually be less then 5 messages a minute.
func Info(msg string, keysAndValues ...interface{}) {
	logger.V(4).Info(msg, keysAndValues...)
}

// Error prints messages about exceptional states of the API or SDK.
func Error(err error, msg string, keysAndValues ...interface{}) {
	logger.Error(err, msg, keysAndValues...)
}

// Debug prints messages about all internal changes in the API or SDK.
func Debug(msg string, keysAndValues ...interface{}) {
	logger.V(8).Info(msg, keysAndValues...)
}

// Warn prints messages about warnings in the API or SDK.
// Not an error but is likely more important than an informational event.
func Warn(msg string, keysAndValues ...interface{}) {
	logger.V(1).Info(msg, keysAndValues...)
}
