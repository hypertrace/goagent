package hypergrpc

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkgrpc "github.com/hypertrace/goagent/sdk/google.golang.org/grpc"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return sdkgrpc.WrapUnaryServerInterceptor(
		otelgrpc.UnaryServerInterceptor(global.Tracer(domain)),
		opentelemetry.SpanFromContext,
	)
}
