package grpcunaryinterceptors

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// InterceptorFilter is a predicate used to determine whether a given request in
// interceptor info should be instrumented. A InterceptorFilter must return true if
// the request should be traced.
//
// Deprecated: Use stats handlers instead.
// type InterceptorFilter func(*otelgrpc.InterceptorInfo) bool

// Filter is a predicate used to determine whether a given request in
// should be instrumented by the attached RPC tag info.
// A Filter must return true if the request should be instrumented.
// type Filter func(*stats.RPCTagInfo) bool

// config is a group of options for this instrumentation.
type config struct {
	Filter            otelgrpc.Filter
	InterceptorFilter otelgrpc.InterceptorFilter
	Propagators       propagation.TextMapPropagator
	TracerProvider    trace.TracerProvider
	MeterProvider     metric.MeterProvider
	SpanStartOptions  []trace.SpanStartOption
	SpanAttributes    []attribute.KeyValue
	MetricAttributes  []attribute.KeyValue

	ReceivedEvent bool
	SentEvent     bool

	tracer trace.Tracer
	meter  metric.Meter

	rpcDuration    metric.Float64Histogram
	rpcInBytes     metric.Int64Histogram
	rpcOutBytes    metric.Int64Histogram
	rpcInMessages  metric.Int64Histogram
	rpcOutMessages metric.Int64Histogram
}

// Option applies an option value for a config.
type Option interface {
	apply(*config)
}

// newConfig returns a config configured with all the passed Options.
func newConfig(opts []Option, role string) *config {
	c := &config{
		Propagators:    otel.GetTextMapPropagator(),
		TracerProvider: otel.GetTracerProvider(),
		MeterProvider:  otel.GetMeterProvider(),
	}
	for _, o := range opts {
		o.apply(c)
	}

	c.tracer = c.TracerProvider.Tracer(
		otelgrpc.ScopeName,
		trace.WithInstrumentationVersion(otelgrpc.Version()),
	)

	c.meter = c.MeterProvider.Meter(
		otelgrpc.ScopeName,
		metric.WithInstrumentationVersion(otelgrpc.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err error
	c.rpcDuration, err = c.meter.Float64Histogram("rpc."+role+".duration",
		metric.WithDescription("Measures the duration of inbound RPC."),
		metric.WithUnit("ms"))
	if err != nil {
		otel.Handle(err)
		if c.rpcDuration == nil {
			c.rpcDuration = noop.Float64Histogram{}
		}
	}

	rpcRequestSize, err := c.meter.Int64Histogram("rpc."+role+".request.size",
		metric.WithDescription("Measures size of RPC request messages (uncompressed)."),
		metric.WithUnit("By"))
	if err != nil {
		otel.Handle(err)
		if rpcRequestSize == nil {
			rpcRequestSize = noop.Int64Histogram{}
		}
	}

	rpcResponseSize, err := c.meter.Int64Histogram("rpc."+role+".response.size",
		metric.WithDescription("Measures size of RPC response messages (uncompressed)."),
		metric.WithUnit("By"))
	if err != nil {
		otel.Handle(err)
		if rpcResponseSize == nil {
			rpcResponseSize = noop.Int64Histogram{}
		}
	}

	rpcRequestsPerRPC, err := c.meter.Int64Histogram("rpc."+role+".requests_per_rpc",
		metric.WithDescription("Measures the number of messages received per RPC. Should be 1 for all non-streaming RPCs."),
		metric.WithUnit("{count}"))
	if err != nil {
		otel.Handle(err)
		if rpcRequestsPerRPC == nil {
			rpcRequestsPerRPC = noop.Int64Histogram{}
		}
	}

	rpcResponsesPerRPC, err := c.meter.Int64Histogram("rpc."+role+".responses_per_rpc",
		metric.WithDescription("Measures the number of messages received per RPC. Should be 1 for all non-streaming RPCs."),
		metric.WithUnit("{count}"))
	if err != nil {
		otel.Handle(err)
		if rpcResponsesPerRPC == nil {
			rpcResponsesPerRPC = noop.Int64Histogram{}
		}
	}

	switch role {
	case "client":
		c.rpcInBytes = rpcResponseSize
		c.rpcInMessages = rpcResponsesPerRPC
		c.rpcOutBytes = rpcRequestSize
		c.rpcOutMessages = rpcRequestsPerRPC
	case "server":
		c.rpcInBytes = rpcRequestSize
		c.rpcInMessages = rpcRequestsPerRPC
		c.rpcOutBytes = rpcResponseSize
		c.rpcOutMessages = rpcResponsesPerRPC
	default:
		c.rpcInBytes = noop.Int64Histogram{}
		c.rpcInMessages = noop.Int64Histogram{}
		c.rpcOutBytes = noop.Int64Histogram{}
		c.rpcOutMessages = noop.Int64Histogram{}
	}

	return c
}
