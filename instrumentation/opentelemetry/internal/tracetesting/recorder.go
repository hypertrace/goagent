package tracetesting

import (
	"context"

	"go.opentelemetry.io/otel/sdk/trace"
)

// Recorder records spans being synced through the SpanSyncer interface.
type Recorder struct {
	spans []trace.ReadOnlySpan
}

var _ trace.SpanExporter = &Recorder{}

func NewRecorder() *Recorder {
	return &Recorder{}
}

// ExportSpans records spans into the internal buffer
func (r *Recorder) ExportSpans(_ context.Context, s []trace.ReadOnlySpan) error {
	r.spans = append(r.spans, s...)
	return nil
}

// Shutdown flushes the buffer
func (r *Recorder) Shutdown(_ context.Context) error {
	_ = r.Flush()
	return nil
}

// Flush returns the current recorded spans and reset the recordings
func (r *Recorder) Flush() []trace.ReadOnlySpan {
	spans := r.spans
	r.spans = nil
	return spans
}
