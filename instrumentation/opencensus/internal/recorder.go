package internal

import "go.opencensus.io/trace"

var _ trace.Exporter = &Recorder{}

// Recorder records spans being exporter through the Exporter interface.
type Recorder struct {
	spans []*trace.SpanData
}

// ExportSpan records a span
func (r *Recorder) ExportSpan(s *trace.SpanData) {
	r.spans = append(r.spans, s)
}

// Flush returns the current recorded spans and reset the recordings
func (r *Recorder) Flush() []*trace.SpanData {
	spans := r.spans
	r.spans = nil
	return spans
}
