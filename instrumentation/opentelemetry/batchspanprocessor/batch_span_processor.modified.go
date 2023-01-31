// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package batchspanprocessor // import "github.com/hypertrace/goagent/instrumentation/opentelemetry/batchspanprocessor"

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	metricglobal "go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Defaults for BatchSpanProcessorOptions.
const (
	DefaultMaxQueueSize       = 2048
	DefaultScheduleDelay      = 5000
	DefaultExportTimeout      = 30000
	DefaultMaxExportBatchSize = 512
	SpansReceivedCounter      = "hypertrace.agent.bsp.spans_received"
	SpansDroppedCounter       = "hypertrace.agent.bsp.spans_dropped"
	SpansUnSampledCounter     = "hypertrace.agent.bsp.spans_unsampled"
)

// batchSpanProcessor is a SpanProcessor that batches asynchronously-received
// spans and sends them to a trace.Exporter when complete.
type batchSpanProcessor struct {
	e sdktrace.SpanExporter
	o sdktrace.BatchSpanProcessorOptions

	queue   chan sdktrace.ReadOnlySpan
	dropped uint32

	batch      []sdktrace.ReadOnlySpan
	batchMutex sync.Mutex
	timer      *time.Timer
	stopWait   sync.WaitGroup
	stopOnce   sync.Once
	stopCh     chan struct{}
	// Some metrics in here.
	counters map[string]syncint64.Counter
}

var _ sdktrace.SpanProcessor = (*batchSpanProcessor)(nil)

// NewBatchSpanProcessor creates a new SpanProcessor that will send completed
// span batches to the exporter with the supplied options.
//
// If the exporter is nil, the span processor will preform no action.
func NewBatchSpanProcessor(exporter sdktrace.SpanExporter, options ...sdktrace.BatchSpanProcessorOption) sdktrace.SpanProcessor {
	maxQueueSize := BatchSpanProcessorMaxQueueSize(DefaultMaxQueueSize)
	maxExportBatchSize := BatchSpanProcessorMaxExportBatchSize(DefaultMaxExportBatchSize)

	if maxExportBatchSize > maxQueueSize {
		if DefaultMaxExportBatchSize > maxQueueSize {
			maxExportBatchSize = maxQueueSize
		} else {
			maxExportBatchSize = DefaultMaxExportBatchSize
		}
	}

	o := sdktrace.BatchSpanProcessorOptions{
		BatchTimeout:       time.Duration(BatchSpanProcessorScheduleDelay(DefaultScheduleDelay)) * time.Millisecond,
		ExportTimeout:      time.Duration(BatchSpanProcessorExportTimeout(DefaultExportTimeout)) * time.Millisecond,
		MaxQueueSize:       maxQueueSize,
		MaxExportBatchSize: maxExportBatchSize,
	}
	for _, opt := range options {
		opt(&o)
	}

	// Setup metrics
	mp := metricglobal.MeterProvider()
	meter := mp.Meter("go.opentelemetry.io/otel/sdk/trace",
		metric.WithInstrumentationVersion(otel.Version()))
	counters := make(map[string]syncint64.Counter)

	// Spans received by processor
	spansReceivedCounter, err := meter.SyncInt64().Counter(SpansReceivedCounter)
	if err != nil {
		otel.Handle(err)
	}
	counters[SpansReceivedCounter] = spansReceivedCounter

	// Spans Dropped by processor once the buffer is full.
	spansDroppedCounter, err := meter.SyncInt64().Counter(SpansDroppedCounter)
	if err != nil {
		otel.Handle(err)
	}
	counters[SpansDroppedCounter] = spansDroppedCounter

	// Spans that are not sampled.(Useful to know when sampling is enabled)
	spansUnSampledCounter, err := meter.SyncInt64().Counter(SpansUnSampledCounter)
	if err != nil {
		otel.Handle(err)
	}
	counters[SpansUnSampledCounter] = spansUnSampledCounter

	bsp := &batchSpanProcessor{
		e:        exporter,
		o:        o,
		batch:    make([]sdktrace.ReadOnlySpan, 0, o.MaxExportBatchSize),
		timer:    time.NewTimer(o.BatchTimeout),
		queue:    make(chan sdktrace.ReadOnlySpan, o.MaxQueueSize),
		stopCh:   make(chan struct{}),
		counters: counters,
	}

	bsp.stopWait.Add(1)
	go func() {
		defer bsp.stopWait.Done()
		bsp.processQueue()
		bsp.drainQueue()
	}()

	return bsp
}

// OnStart method does nothing.
func (bsp *batchSpanProcessor) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {}

// OnEnd method enqueues a ReadOnlySpan for later processing.
func (bsp *batchSpanProcessor) OnEnd(s sdktrace.ReadOnlySpan) {
	// Do not enqueue spans if we are just going to drop them.
	if bsp.e == nil {
		return
	}
	bsp.enqueue(s)
}

// Shutdown flushes the queue and waits until all spans are processed.
// It only executes once. Subsequent call does nothing.
func (bsp *batchSpanProcessor) Shutdown(ctx context.Context) error {
	var err error
	bsp.stopOnce.Do(func() {
		wait := make(chan struct{})
		go func() {
			close(bsp.stopCh)
			bsp.stopWait.Wait()
			if bsp.e != nil {
				if err := bsp.e.Shutdown(ctx); err != nil {
					otel.Handle(err)
				}
			}
			close(wait)
		}()
		// Wait until the wait group is done or the context is cancelled
		select {
		case <-wait:
		case <-ctx.Done():
			err = ctx.Err()
		}
	})
	return err
}

type forceFlushSpan struct {
	sdktrace.ReadOnlySpan
	flushed chan struct{}
}

func (f forceFlushSpan) SpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{TraceFlags: trace.FlagsSampled})
}

// ForceFlush exports all ended spans that have not yet been exported.
func (bsp *batchSpanProcessor) ForceFlush(ctx context.Context) error {
	var err error
	if bsp.e != nil {
		flushCh := make(chan struct{})
		if bsp.enqueueBlockOnQueueFull(ctx, forceFlushSpan{flushed: flushCh}) {
			select {
			case <-flushCh:
				// Processed any items in queue prior to ForceFlush being called
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		wait := make(chan error)
		go func() {
			wait <- bsp.exportSpans(ctx)
			close(wait)
		}()
		// Wait until the export is finished or the context is cancelled/timed out
		select {
		case err = <-wait:
		case <-ctx.Done():
			err = ctx.Err()
		}
	}
	return err
}

// WithMaxQueueSize returns a BatchSpanProcessorOption that configures the
// maximum queue size allowed for a BatchSpanProcessor.
func WithMaxQueueSize(size int) sdktrace.BatchSpanProcessorOption {
	return func(o *sdktrace.BatchSpanProcessorOptions) {
		o.MaxQueueSize = size
	}
}

// WithMaxExportBatchSize returns a BatchSpanProcessorOption that configures
// the maximum export batch size allowed for a BatchSpanProcessor.
func WithMaxExportBatchSize(size int) sdktrace.BatchSpanProcessorOption {
	return func(o *sdktrace.BatchSpanProcessorOptions) {
		o.MaxExportBatchSize = size
	}
}

// WithBatchTimeout returns a BatchSpanProcessorOption that configures the
// maximum delay allowed for a BatchSpanProcessor before it will export any
// held span (whether the queue is full or not).
func WithBatchTimeout(delay time.Duration) sdktrace.BatchSpanProcessorOption {
	return func(o *sdktrace.BatchSpanProcessorOptions) {
		o.BatchTimeout = delay
	}
}

// WithExportTimeout returns a BatchSpanProcessorOption that configures the
// amount of time a BatchSpanProcessor waits for an exporter to export before
// abandoning the export.
func WithExportTimeout(timeout time.Duration) sdktrace.BatchSpanProcessorOption {
	return func(o *sdktrace.BatchSpanProcessorOptions) {
		o.ExportTimeout = timeout
	}
}

// WithBlocking returns a BatchSpanProcessorOption that configures a
// BatchSpanProcessor to wait for enqueue operations to succeed instead of
// dropping data when the queue is full.
func WithBlocking() sdktrace.BatchSpanProcessorOption {
	return func(o *sdktrace.BatchSpanProcessorOptions) {
		o.BlockOnQueueFull = true
	}
}

// exportSpans is a subroutine of processing and draining the queue.
func (bsp *batchSpanProcessor) exportSpans(ctx context.Context) error {
	bsp.timer.Reset(bsp.o.BatchTimeout)

	bsp.batchMutex.Lock()
	defer bsp.batchMutex.Unlock()

	if bsp.o.ExportTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, bsp.o.ExportTimeout)
		defer cancel()
	}

	if l := len(bsp.batch); l > 0 {
		Debug("exporting spans", "count", len(bsp.batch), "total_dropped", atomic.LoadUint32(&bsp.dropped))
		err := bsp.e.ExportSpans(ctx, bsp.batch)

		// A new batch is always created after exporting, even if the batch failed to be exported.
		//
		// It is up to the exporter to implement any type of retry logic if a batch is failing
		// to be exported, since it is specific to the protocol and backend being sent to.
		bsp.batch = bsp.batch[:0]

		if err != nil {
			return err
		}
	}
	return nil
}

// processQueue removes spans from the `queue` channel until processor
// is shut down. It calls the exporter in batches of up to MaxExportBatchSize
// waiting up to BatchTimeout to form a batch.
func (bsp *batchSpanProcessor) processQueue() {
	defer bsp.timer.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case <-bsp.stopCh:
			return
		case <-bsp.timer.C:
			if err := bsp.exportSpans(ctx); err != nil {
				otel.Handle(err)
			}
		case sd := <-bsp.queue:
			if ffs, ok := sd.(forceFlushSpan); ok {
				close(ffs.flushed)
				continue
			}
			bsp.batchMutex.Lock()
			bsp.batch = append(bsp.batch, sd)
			shouldExport := len(bsp.batch) >= bsp.o.MaxExportBatchSize
			bsp.batchMutex.Unlock()
			if shouldExport {
				if !bsp.timer.Stop() {
					<-bsp.timer.C
				}
				if err := bsp.exportSpans(ctx); err != nil {
					otel.Handle(err)
				}
			}
		}
	}
}

// drainQueue awaits the any caller that had added to bsp.stopWait
// to finish the enqueue, then exports the final batch.
func (bsp *batchSpanProcessor) drainQueue() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case sd := <-bsp.queue:
			if sd == nil {
				if err := bsp.exportSpans(ctx); err != nil {
					otel.Handle(err)
				}
				return
			}

			bsp.batchMutex.Lock()
			bsp.batch = append(bsp.batch, sd)
			shouldExport := len(bsp.batch) == bsp.o.MaxExportBatchSize
			bsp.batchMutex.Unlock()

			if shouldExport {
				if err := bsp.exportSpans(ctx); err != nil {
					otel.Handle(err)
				}
			}
		default:
			close(bsp.queue)
		}
	}
}

func (bsp *batchSpanProcessor) enqueue(sd sdktrace.ReadOnlySpan) {
	ctx := context.TODO()
	if bsp.o.BlockOnQueueFull {
		bsp.enqueueBlockOnQueueFull(ctx, sd)
	} else {
		bsp.enqueueDrop(ctx, sd)
	}
}

func recoverSendOnClosedChan() {
	x := recover()
	switch err := x.(type) {
	case nil:
		return
	case runtime.Error:
		if err.Error() == "send on closed channel" {
			return
		}
	}
	panic(x)
}

func (bsp *batchSpanProcessor) enqueueBlockOnQueueFull(ctx context.Context, sd sdktrace.ReadOnlySpan) bool {
	if !sd.SpanContext().IsSampled() {
		return false
	}

	// This ensures the bsp.queue<- below does not panic as the
	// processor shuts down.
	defer recoverSendOnClosedChan()

	select {
	case <-bsp.stopCh:
		return false
	default:
	}

	select {
	case bsp.queue <- sd:
		return true
	case <-ctx.Done():
		return false
	}
}

func (bsp *batchSpanProcessor) enqueueDrop(ctx context.Context, sd sdktrace.ReadOnlySpan) bool {
	if !sd.SpanContext().IsSampled() {
		// Count the span as unsampled
		bsp.counters[SpansUnSampledCounter].Add(ctx, 1)
		return false
	}

	// Count the span as received.
	bsp.counters[SpansReceivedCounter].Add(ctx, 1)

	// This ensures the bsp.queue<- below does not panic as the
	// processor shuts down.
	defer recoverSendOnClosedChan()

	select {
	case <-bsp.stopCh:
		return false
	default:
	}

	select {
	case bsp.queue <- sd:
		return true
	default:
		atomic.AddUint32(&bsp.dropped, 1)
		// Count the span as dropped.
		bsp.counters[SpansDroppedCounter].Add(ctx, 1)
	}
	return false
}

// MarshalLog is the marshaling function used by the logging system to represent this exporter.
func (bsp *batchSpanProcessor) MarshalLog() interface{} {
	return struct {
		Type         string
		SpanExporter sdktrace.SpanExporter
		Config       sdktrace.BatchSpanProcessorOptions
	}{
		Type:         "BatchSpanProcessor",
		SpanExporter: bsp.e,
		Config:       bsp.o,
	}
}
