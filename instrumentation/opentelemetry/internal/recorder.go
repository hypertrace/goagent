package internal

import (
	"context"

	"go.opentelemetry.io/otel/sdk/export/trace"
)

var _ trace.SpanSyncer = &Recorder{}

// Recorder records spans being synced through the SpanSyncer interface.
type Recorder struct {
	spans []*trace.SpanData
}

// ExportSpan records a span
func (r *Recorder) ExportSpan(ctx context.Context, s *trace.SpanData) {
	r.spans = append(r.spans, s)
}

// Flush returns the current recorded spans and reset the recordings
func (r *Recorder) Flush() []*trace.SpanData {
	spans := r.spans
	r.spans = nil
	return spans
}
