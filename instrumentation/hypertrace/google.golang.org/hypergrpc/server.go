package hypergrpc // import "github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc"

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkgrpc "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return sdkgrpc.WrapUnaryServerInterceptor(
		otelgrpc.UnaryServerInterceptor(),
		opentelemetry.SpanFromContext,
		o.toSDKOptions(),
		map[string]string{},
	)
}
