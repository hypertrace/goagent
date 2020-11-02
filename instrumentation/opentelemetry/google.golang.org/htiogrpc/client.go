package htiogrpc

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkgrpc "github.com/hypertrace/goagent/sdk/google.golang.org/grpc"
	"google.golang.org/grpc"
)

// WrapUnaryClientInterceptor returns a new unary client interceptor that will
// complement existing OpenTelemetry instrumentation
func WrapUnaryClientInterceptor(delegate grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return sdkgrpc.WrapUnaryClientInterceptor(delegate, opentelemetry.SpanFromContext)
}
