# Modified Batch Span Processor

Since the original BatchSpanProcessor does not send metrics for spans received and spans dropped, the modified BatchSpanProcessor creates and populates those metrics so the user can track whether there are spans being dropped.

We have kept track of the original files modified so it's easier to figure out the changes we added. When upgrading, copy the newer files into their go.original counterparts and then do a diff with the modified.go files to figure out what changes to make in order to upgrade the modified files.

The paths of the files modified:
- [sdk/trace/batch_span_processor.go](https://github.com/open-telemetry/opentelemetry-go/blob/main/sdk/trace/batch_span_processor.go)
- [sdk/internal/env/env.go](https://github.com/open-telemetry/opentelemetry-go/blob/main/sdk/internal/env/env.go)

Since we cannot use [the internal logger]((https://github.com/open-telemetry/opentelemetry-go/blob/main/internal/global/internal_logging.go)), we have adapted it at [logger.go](instrumentation/opentelemetry/batchspanprocessor/logger.go).
