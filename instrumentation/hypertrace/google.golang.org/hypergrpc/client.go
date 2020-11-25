package hypergrpc

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkgrpc "github.com/hypertrace/goagent/sdk/google.golang.org/grpc"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor suitable
// for use in a grpc.Dial call.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return sdkgrpc.WrapUnaryClientInterceptor(
		otelgrpc.UnaryClientInterceptor(global.Tracer(opentelemetry.TracerDomain)),
		opentelemetry.SpanFromContext,
	)
}
