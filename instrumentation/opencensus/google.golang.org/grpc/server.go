package grpc

import (
	"github.com/traceableai/goagent/instrumentation/opencensus"
	traceablegrpc "github.com/traceableai/goagent/sdk/google.golang.org/grpc"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptor that will
// complement existing OpenCensus instrumentation
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return traceablegrpc.EnrichUnaryServerInterceptor(nil, opencensus.SpanFromContext)
}
