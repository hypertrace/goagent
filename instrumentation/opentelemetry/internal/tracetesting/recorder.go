package tracetesting

import (
	"context"

	"go.opentelemetry.io/otel/sdk/trace"
)

var _ trace.SpanExporter = &recorder{}

// recorder records spans being synced through the SpanSyncer interface.
type recorder struct {
	spans []trace.ReadOnlySpan
}

// ExportSpans records spans into the internal buffer
func (r *recorder) ExportSpans(_ context.Context, s []trace.ReadOnlySpan) error {
	r.spans = append(r.spans, s...)
	return nil
}

// Shutdown flushes the buffer
func (r *recorder) Shutdown(_ context.Context) error {
	_ = r.Flush()
	return nil
}

// Flush returns the current recorded spans and reset the recordings
func (r *recorder) Flush() []trace.ReadOnlySpan {
	spans := r.spans
	r.spans = nil
	return spans
}
