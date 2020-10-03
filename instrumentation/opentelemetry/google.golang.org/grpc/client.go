package grpc

import (
	"github.com/traceableai/goagent/instrumentation/opentelemetry"
	traceablegrpc "github.com/traceableai/goagent/sdk/google.golang.org/grpc"
	"google.golang.org/grpc"
)

// EnrichUnaryClientInterceptor returns a new unary client interceptor that will
// complement existing OpenTelemetry instrumentation
func EnrichUnaryClientInterceptor(delegate grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return traceablegrpc.EnrichUnaryClientInterceptor(delegate, opentelemetry.SpanFromContext)
}
