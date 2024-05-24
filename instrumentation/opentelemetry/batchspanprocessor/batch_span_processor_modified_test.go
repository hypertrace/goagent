package batchspanprocessor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	panicSpanStr  string = "panic_span"
	tracerNameStr string = "tracer1"
)

func TestCustomBspNonPanicExporterShouldNotPanic(t *testing.T) {
	tracer, verifyFunc := setupTracer(t, true, false)

	startAndEndSpan(tracer, "span1")
	startAndEndSpan(tracer, "span2")
	startAndEndSpan(tracer, panicSpanStr)
	time.Sleep(8 * time.Millisecond)

	verifyFunc(3)
	startAndEndSpan(tracer, "span4")
	time.Sleep(8 * time.Millisecond)
	verifyFunc(4)
}

func TestCustomBspPanicExporterShouldNotPanic(t *testing.T) {
	tracer, verifyFunc := setupTracer(t, true, true)

	startAndEndSpan(tracer, "span1")
	startAndEndSpan(tracer, "span2")
	startAndEndSpan(tracer, panicSpanStr)
	time.Sleep(8 * time.Millisecond)
	verifyFunc(2)

	startAndEndSpan(tracer, "span4")
	time.Sleep(8 * time.Millisecond)
	verifyFunc(3)

	// Only span5 and span6 will be exported. span8 is discarded since the spans loop does not get to it before
	// the panic
	startAndEndSpan(tracer, "span5")
	startAndEndSpan(tracer, "span6")
	startAndEndSpan(tracer, panicSpanStr)
	startAndEndSpan(tracer, "span8")
	time.Sleep(8 * time.Millisecond)
	verifyFunc(5)

	startAndEndSpan(tracer, "span9")
	startAndEndSpan(tracer, "span10")
	startAndEndSpan(tracer, "span11")
	time.Sleep(8 * time.Millisecond)
	verifyFunc(8)
}

func TestCustomBspPanicExporterGoodSpansShouldNotPanic(t *testing.T) {
	tracer, verifyFunc := setupTracer(t, true, false)

	startAndEndSpan(tracer, "span1")
	startAndEndSpan(tracer, "span2")
	startAndEndSpan(tracer, "span3")
	time.Sleep(8 * time.Millisecond)

	verifyFunc(3)

	startAndEndSpan(tracer, "span4")
	time.Sleep(8 * time.Millisecond)
	verifyFunc(4)
}

func setupTracer(t *testing.T, useCustomBsp bool, enablePanic bool) (trace.Tracer, func(int)) {
	exporter := &mockPanickingSpanExporter{panics: enablePanic}
	exportTimeout := 5 * time.Millisecond
	bsp := CreateBatchSpanProcessor(useCustomBsp, exporter,
		sdktrace.WithBatchTimeout(exportTimeout))
	assert.NotNil(t, bsp)

	sampler := sdktrace.AlwaysSample()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithSpanProcessor(bsp),
	)

	return tp.Tracer(tracerNameStr), func(expectedExportedSpans int) {
		verifyExporter(t, exporter, expectedExportedSpans)
	}
}

func verifyExporter(t *testing.T, e *mockPanickingSpanExporter, expectedSpansExporter int) {
	assert.Equal(t, expectedSpansExporter, e.exportedCount)
}

func startAndEndSpan(tracer trace.Tracer, spanName string) {
	_, s := tracer.Start(context.Background(), spanName)
	time.Sleep(1 * time.Millisecond)
	s.End()
}

type mockPanickingSpanExporter struct {
	panics        bool
	exportedCount int
}

func (e *mockPanickingSpanExporter) ExportSpans(_ context.Context, spans []sdktrace.ReadOnlySpan) error {
	if !e.panics {
		e.exportedCount = e.exportedCount + len(spans)
		return nil
	}

	for _, span := range spans {
		if span.Name() == panicSpanStr {
			panic("panic span in span list")
		} else {
			e.exportedCount = e.exportedCount + 1
		}
	}
	return nil
}

func (e *mockPanickingSpanExporter) Shutdown(_ context.Context) error {
	return nil
}
