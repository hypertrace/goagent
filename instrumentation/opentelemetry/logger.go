package opentelemetry

import (
	"context"
	"io"
	//"log"
	//"os"
	"time"

	//"github.com/go-logr/logr"
	//"github.com/go-logr/logr/funcr"
	//"github.com/go-logr/stdr"
	"github.com/hypertrace/goagent/sdk"
)

//var globalLogger logr.Logger = stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))

// type LogSink interface {
// 	// Init receives optional information about the logr library for LogSink
// 	// implementations that need it.
// 	Init(info logr.RuntimeInfo)

// 	// Enabled tests whether this LogSink is enabled at the specified V-level.
// 	// For example, commandline flags might be used to set the logging
// 	// verbosity and disable some info logs.
// 	Enabled(level int) bool

// 	// Info logs a non-error message with the given key/value pairs as context.
// 	// The level argument is provided for optional logging.  This method will
// 	// only be called when Enabled(level) is true. See Logger.Info for more
// 	// details.
// 	Info(level int, msg string, keysAndValues ...interface{})

// 	// Error logs an error, with the given message and key/value pairs as
// 	// context.  See Logger.Error for more details.
// 	Error(err error, msg string, keysAndValues ...interface{})

// 	// WithValues returns a new LogSink with additional key/value pairs.  See
// 	// Logger.WithValues for more details.
// 	WithValues(keysAndValues ...interface{}) LogSink

// 	// WithName returns a new LogSink with the specified name appended.  See
// 	// Logger.WithName for more details.
// 	WithName(name string) LogSink
// }

type loggerSink struct {
	//funcr.Formatter
	//std            stdr.StdLogger
	logEntriesChan chan string
	delegate       io.Writer
}

func newLoggerSink(delegate io.Writer) *loggerSink {
	//stdrLogger := stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))
	return &loggerSink{make(chan string, 2048), delegate}
}

func (l *loggerSink) Write(p []byte) (n int, err error) {
	pLen, err := l.delegate.Write(p)
	if err != nil {
		return pLen, err
	}

	pLen = len(p)
	if pLen == 0 {
		return 0, nil
	}
	pCopy := make([]byte, pLen)
	copy(pCopy, p)
	l.logEntriesChan <- string(pCopy)
	return pLen, nil
}

// func (l *LogsIntoSpansWriter) Sync() error {
// 	return nil
// }

// LogsAsSpans converts traceable-agent's logs to spans which are reported to the platform.
func (l *loggerSink) LogsAsSpans(startSpan sdk.StartSpan) error {
	ticker := time.NewTicker(5 * time.Second)
	spanPresent := false
	var ender func()
	var span sdk.Span
	//opts := []goagent.Option{}
	for {
		select {
		case le := <-l.logEntriesChan:
			if !spanPresent {
				_, span, ender = startSpan(context.Background(), "log",
					&sdk.SpanOptions{Kind: sdk.SpanKindUndetermined, Timestamp: time.Now()})
				span.SetAttribute("deployment.environment", "traceableai-internal")
				spanPresent = true
			}
			span.AddEvent(le, time.Now(), map[string]interface{}{})
		case <-ticker.C:
			if spanPresent {
				ender()
				spanPresent = false
			}
		}
	}
}

// func (l *LogsIntoSpansWriter) GetLoggingOptions(config *configv1.TraceableServiceGlobalConfig) []zap.Option {
// 	if !l.Enabled {
// 		return []zap.Option{}
// 	}

// 	zapConfig, err := createZapConfig(config)
// 	if err != nil {
// 		return []zap.Option{}
// 	}

// 	return []zap.Option{multiWriteSyncerOpt(l, zapConfig)}
// }

// func (l *LogsIntoSpansWriter) GetLoggerWithLevel(config *configv1.TraceableServiceGlobalConfig, level zapcore.Level) (*zap.Logger, error) {
// 	return createLogger(config, level, l, true)
// }

// func multiWriteSyncerOpt(logsIntoSpansWriter *LogsIntoSpansWriter, zapConfig zap.Config) zap.Option {
// 	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
// 		return zapcore.NewCore(
// 			zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
// 			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logsIntoSpansWriter),
// 			zapConfig.Level,
// 		)
// 	})
// }
