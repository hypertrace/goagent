package opentelemetry // import "github.com/hypertrace/goagent/instrumentation/opentelemetry"

import (
	"context"
	"time"

	"github.com/hypertrace/goagent/sdk"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ sdk.Span = (*Span)(nil)

type Span struct {
	trace.Span
}

func (s *Span) SetAttribute(key string, value interface{}) {
	s.Span.SetAttributes(attribute.Any(key, value))
}

func (s *Span) SetError(err error) {
	s.Span.RecordError(err)
}

func (s *Span) SetStatus(code sdk.Code, description string) {
	s.Span.SetStatus(codes.Code(code), description)
}

func (s *Span) IsNoop() bool {
	return !s.Span.IsRecording()
}

func SpanFromContext(ctx context.Context) sdk.Span {
	return &Span{trace.SpanFromContext(ctx)}
}

type getTracerProvider func() trace.TracerProvider

func startSpan(provider getTracerProvider) sdk.StartSpan {
	return func(ctx context.Context, name string, opts *sdk.SpanOptions) (context.Context, sdk.Span, func()) {
		startOpts := []trace.SpanStartOption{}

		if opts != nil {
			startOpts = append(startOpts, trace.WithSpanKind(mapSpanKind(opts.Kind)))

			if opts.Timestamp.IsZero() {
				startOpts = append(startOpts, trace.WithTimestamp(time.Now()))
			} else {
				startOpts = append(startOpts, trace.WithTimestamp(opts.Timestamp))
			}
		}

		ctx, span := provider().Tracer(TracerDomain).Start(ctx, name, startOpts...)
		return ctx, &Span{span}, func() { span.End() }
	}
}

var StartSpan = startSpan(otel.GetTracerProvider)

func mapSpanKind(kind sdk.SpanKind) trace.SpanKind {
	switch kind {
	case sdk.SpanKindClient:
		return trace.SpanKindClient
	case sdk.SpanKindServer:
		return trace.SpanKindServer
	case sdk.SpanKindProducer:
		return trace.SpanKindProducer
	case sdk.SpanKindConsumer:
		return trace.SpanKindConsumer
	default:
		return trace.SpanKindUnspecified
	}
}
