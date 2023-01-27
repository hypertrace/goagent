package batchspanprocessor

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func CreateBatchSpanProcessor(useModified bool, exporter sdktrace.SpanExporter,
	options ...sdktrace.BatchSpanProcessorOption) sdktrace.SpanProcessor {
	if useModified {
		return NewBatchSpanProcessor(exporter, options...)
	} else {
		return sdktrace.NewBatchSpanProcessor(exporter, options...)
	}
}
