package grpcunaryinterceptors

// Copy two functions from https://github.com/open-telemetry/opentelemetry-go-contrib/blob/v1.34.0/instrumentation/google.golang.org/grpc/otelgrpc/interceptor.go
// UnaryClientInterceptor and UnaryServerInterceptor were deleted. We rely on them for instrumentation of grpc so we will copy the deleted functions here
// temporarily while we switch to the stats.Handler interface.
//
// This package should not be here for long. Once we switch to the stats.Handler interface, we can remove this package.

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type messageType attribute.KeyValue

var (
	messageSent     = messageType(RPCMessageTypeSent)
	messageReceived = messageType(RPCMessageTypeReceived)
)

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor suitable
// for use in a grpc.NewClient call.
//
// Deprecated: Use [NewClientHandler] instead.
func UnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
	cfg := newConfig(opts, "client")
	tracer := cfg.TracerProvider.Tracer(
		otelgrpc.ScopeName,
		trace.WithInstrumentationVersion(otelgrpc.Version()),
	)

	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		callOpts ...grpc.CallOption,
	) error {
		i := &otelgrpc.InterceptorInfo{
			Method: method,
			Type:   otelgrpc.UnaryClient,
		}
		if cfg.InterceptorFilter != nil && !cfg.InterceptorFilter(i) {
			return invoker(ctx, method, req, reply, cc, callOpts...)
		}

		name, attr, _ := telemetryAttributes(method, cc.Target())

		startOpts := append([]trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attr...),
		},
			cfg.SpanStartOptions...,
		)

		ctx, span := tracer.Start(
			ctx,
			name,
			startOpts...,
		)
		defer span.End()

		ctx = inject(ctx, cfg.Propagators)

		if cfg.SentEvent {
			messageSent.Event(ctx, 1, req)
		}

		err := invoker(ctx, method, req, reply, cc, callOpts...)

		if cfg.ReceivedEvent {
			messageReceived.Event(ctx, 1, reply)
		}

		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(statusCodeAttr(s.Code()))
		} else {
			span.SetAttributes(statusCodeAttr(grpc_codes.OK))
		}

		return err
	}
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
//
// Deprecated: Use [NewServerHandler] instead.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	cfg := newConfig(opts, "server")
	tracer := cfg.TracerProvider.Tracer(
		otelgrpc.ScopeName,
		trace.WithInstrumentationVersion(otelgrpc.Version()),
	)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		i := &otelgrpc.InterceptorInfo{
			UnaryServerInfo: info,
			Type:            otelgrpc.UnaryServer,
		}
		if cfg.InterceptorFilter != nil && !cfg.InterceptorFilter(i) {
			return handler(ctx, req)
		}

		ctx = extract(ctx, cfg.Propagators)
		name, attr, metricAttrs := telemetryAttributes(info.FullMethod, peerFromCtx(ctx))

		startOpts := append([]trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attr...),
		},
			cfg.SpanStartOptions...,
		)

		ctx, span := tracer.Start(
			trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
			name,
			startOpts...,
		)
		defer span.End()

		if cfg.ReceivedEvent {
			messageReceived.Event(ctx, 1, req)
		}

		before := time.Now()

		resp, err := handler(ctx, req)

		s, _ := status.FromError(err)
		if err != nil {
			statusCode, msg := serverStatus(s)
			span.SetStatus(statusCode, msg)
			if cfg.SentEvent {
				messageSent.Event(ctx, 1, s.Proto())
			}
		} else {
			if cfg.SentEvent {
				messageSent.Event(ctx, 1, resp)
			}
		}
		grpcStatusCodeAttr := statusCodeAttr(s.Code())
		span.SetAttributes(grpcStatusCodeAttr)

		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedTime := float64(time.Since(before)) / float64(time.Millisecond)

		metricAttrs = append(metricAttrs, grpcStatusCodeAttr)
		cfg.rpcDuration.Record(ctx, elapsedTime, metric.WithAttributeSet(attribute.NewSet(metricAttrs...)))

		return resp, err
	}
}

// Event adds an event of the messageType to the span associated with the
// passed context with a message id.
func (m messageType) Event(ctx context.Context, id int, _ interface{}) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	span.AddEvent("message", trace.WithAttributes(
		attribute.KeyValue(m),
		RPCMessageIDKey.Int(id),
	))
}

// telemetryAttributes returns a span name and span and metric attributes from
// the gRPC method and peer address.
func telemetryAttributes(fullMethod, peerAddress string) (string, []attribute.KeyValue, []attribute.KeyValue) {
	name, methodAttrs := internalParseFullMethod(fullMethod)
	peerAttrs := peerAttr(peerAddress)

	attrs := make([]attribute.KeyValue, 0, 1+len(methodAttrs)+len(peerAttrs))
	attrs = append(attrs, RPCSystemGRPC)
	attrs = append(attrs, methodAttrs...)
	metricAttrs := attrs[:1+len(methodAttrs)]
	attrs = append(attrs, peerAttrs...)
	return name, attrs, metricAttrs
}

// peerAttr returns attributes about the peer address.
func peerAttr(addr string) []attribute.KeyValue {
	host, p, err := net.SplitHostPort(addr)
	if err != nil {
		return nil
	}

	if host == "" {
		host = "127.0.0.1"
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return nil
	}

	var attr []attribute.KeyValue
	if ip := net.ParseIP(host); ip != nil {
		attr = []attribute.KeyValue{
			semconv.NetSockPeerAddr(host),
			semconv.NetSockPeerPort(port),
		}
	} else {
		attr = []attribute.KeyValue{
			semconv.NetPeerName(host),
			semconv.NetPeerPort(port),
		}
	}

	return attr
}

// peerFromCtx returns a peer address from a context, if one exists.
func peerFromCtx(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	return p.Addr.String()
}

// [Removed in 1.36.0]
// GRPCStatusCodeKey is convention for numeric status code of a gRPC request.
const GRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")

// statusCodeAttr returns status code attribute based on given gRPC code.
func statusCodeAttr(c grpc_codes.Code) attribute.KeyValue {
	return GRPCStatusCodeKey.Int64(int64(c))
}

// serverStatus returns a span status code and message for a given gRPC
// status code. It maps specific gRPC status codes to a corresponding span
// status code and message. This function is intended for use on the server
// side of a gRPC connection.
//
// If the gRPC status code is Unknown, DeadlineExceeded, Unimplemented,
// Internal, Unavailable, or DataLoss, it returns a span status code of Error
// and the message from the gRPC status. Otherwise, it returns a span status
// code of Unset and an empty message.
func serverStatus(grpcStatus *status.Status) (codes.Code, string) {
	switch grpcStatus.Code() {
	case grpc_codes.Unknown,
		grpc_codes.DeadlineExceeded,
		grpc_codes.Unimplemented,
		grpc_codes.Internal,
		grpc_codes.Unavailable,
		grpc_codes.DataLoss:
		return codes.Error, grpcStatus.Message()
	default:
		return codes.Unset, ""
	}
}

// ParseFullMethod returns a span name following the OpenTelemetry semantic
// conventions as well as all applicable span attribute.KeyValue attributes based
// on a gRPC's FullMethod.
//
// Parsing is consistent with grpc-go implementation:
// https://github.com/grpc/grpc-go/blob/v1.57.0/internal/grpcutil/method.go#L26-L39
//
// Copied from go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/internal#ParseFullMethod
func internalParseFullMethod(fullMethod string) (string, []attribute.KeyValue) {
	if !strings.HasPrefix(fullMethod, "/") {
		// Invalid format, does not follow `/package.service/method`.
		return fullMethod, nil
	}
	name := fullMethod[1:]
	pos := strings.LastIndex(name, "/")
	if pos < 0 {
		// Invalid format, does not follow `/package.service/method`.
		return name, nil
	}
	service, method := name[:pos], name[pos+1:]

	var attrs []attribute.KeyValue
	if service != "" {
		attrs = append(attrs, semconv.RPCService(service))
	}
	if method != "" {
		attrs = append(attrs, semconv.RPCMethod(method))
	}
	return name, attrs
}
