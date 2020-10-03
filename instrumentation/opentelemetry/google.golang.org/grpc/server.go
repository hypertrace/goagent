package grpc

import (
	"github.com/traceableai/goagent/instrumentation/opentelemetry"
	traceablegrpc "github.com/traceableai/goagent/sdk/google.golang.org/grpc"
	"google.golang.org/grpc"
)

// EnrichUnaryServerInterceptor returns a new unary server interceptor tthat will
// complement existing OpenTelemetry instrumentation
func EnrichUnaryServerInterceptor(delegate grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return traceablegrpc.EnrichUnaryServerInterceptor(delegate, opentelemetry.SpanFromContext)
}
