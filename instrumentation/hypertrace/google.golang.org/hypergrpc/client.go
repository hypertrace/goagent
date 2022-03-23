package hypergrpc // import "github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc"

import (
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	sdkgrpc "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor suitable
// for use in a grpc.Dial call.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return sdkgrpc.WrapUnaryClientInterceptor(
		otelgrpc.UnaryClientInterceptor(),
		opentelemetry.SpanFromContext,
		map[string]string{},
	)
}
