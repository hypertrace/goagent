package internal

import "go.opencensus.io/trace"

// InitTracer initializes the tracer and returns a flusher of the reported
// spans for further inspection. Its main purpose is to declare a tracer
// for TESTING.
func InitTracer() func() []*trace.SpanData {
	recorder := &Recorder{}
	trace.RegisterExporter(recorder)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return recorder.Flush
}
