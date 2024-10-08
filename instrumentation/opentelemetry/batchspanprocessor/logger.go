package batchspanprocessor // import "github.com/hypertrace/goagent/instrumentation/opentelemetry/batchspanprocessor"

// Adapted from go.opentelemetry.io/otel/internal/global#internal_logging.go
import (
	"sync"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

// The logger uses stdr which is backed by the standard `log.Logger`
// interface. This logger will only show messages at the Error Level.
var (
	logger logr.Logger
	once   sync.Once
)

// Info prints messages about the general state of the API or SDK.
// This should usually be less then 5 messages a minute.
func Info(msg string, keysAndValues ...interface{}) {
	getLogger().V(4).Info(msg, keysAndValues...)
}

// Error prints messages about exceptional states of the API or SDK.
func Error(err error, msg string, keysAndValues ...interface{}) {
	getLogger().Error(err, msg, keysAndValues...)
}

// Debug prints messages about all internal changes in the API or SDK.
func Debug(msg string, keysAndValues ...interface{}) {
	getLogger().V(8).Info(msg, keysAndValues...)
}

// Warn prints messages about warnings in the API or SDK.
// Not an error but is likely more important than an informational event.
func Warn(msg string, keysAndValues ...interface{}) {
	getLogger().V(1).Info(msg, keysAndValues...)
}

// since global variables init happens on package imports,
// the zap global logger might not be set by then
// hence doing a lazy initialization
func getLogger() logr.Logger {
	once.Do(func() {
		logger = zapr.NewLogger(zap.L())
	})
	return logger
}
