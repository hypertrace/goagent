package hypergrpc // import "github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc"

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/grpcunaryinterceptors"
	sdkgrpc "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor suitable
// for use in a grpc.Dial call.
// Interceptor format will be replaced with the stats.Handler since instrumentation has moved to the stats.Handler.
// See: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/v1.36.0/instrumentation/google.golang.org/grpc/otelgrpc/example_test.go
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return sdkgrpc.WrapUnaryClientInterceptor(
		grpcunaryinterceptors.UnaryClientInterceptor(),
		opentelemetry.SpanFromContext,
		map[string]string{},
	)
}
