package grpc

import (
	otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/grpc"
)

// NewUnaryClientInterceptor returns a new unary client interceptor
func NewUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return otel.UnaryClientInterceptor(
		global.TraceProvider().Tracer("ai.traceable"),
	)
}

// NewStreamClientInterceptor returns a new stream client interceptor
func NewStreamClientInterceptor() grpc.StreamClientInterceptor {
	return otel.StreamClientInterceptor(
		global.TraceProvider().Tracer("ai.traceable"),
	)
}

// NewUnaryServerInterceptor returns a new unary server interceptor
func NewUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return otel.UnaryServerInterceptor(
		global.TraceProvider().Tracer("ai.traceable"),
	)
}

// NewStreamServerInterceptor returns a new stream server interceptor
func NewStreamServerInterceptor() grpc.StreamServerInterceptor {
	return otel.StreamServerInterceptor(
		global.TraceProvider().Tracer("ai.traceable"),
	)
}
