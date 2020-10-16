package grpc

import (
	"github.com/traceableai/goagent/instrumentation/opentelemetry"
	traceablegrpc "github.com/traceableai/goagent/sdk/google.golang.org/grpc"
	"google.golang.org/grpc"
)

// WrapUnaryServerInterceptor returns a new unary server interceptor that will
// complement existing OpenTelemetry instrumentation
func WrapUnaryServerInterceptor(delegate grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return traceablegrpc.WrapUnaryServerInterceptor(delegate, opentelemetry.SpanFromContext)
}
