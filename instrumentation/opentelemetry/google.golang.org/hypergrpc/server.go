package hypergrpc

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkgrpc "github.com/hypertrace/goagent/sdk/google.golang.org/grpc"
	"google.golang.org/grpc"
)

// WrapUnaryServerInterceptor returns a new unary server interceptor that will
// complement existing OpenTelemetry instrumentation
func WrapUnaryServerInterceptor(delegate grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return sdkgrpc.WrapUnaryServerInterceptor(delegate, opentelemetry.SpanFromContext)
}
