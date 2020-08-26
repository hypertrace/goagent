package client

import (
	"context"

	"github.com/traceableai/goagent/otel/grpc/internal"
	otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc"
)

// NewUnaryClientInterceptor returns an interceptor that records the request and response message's body
// and serialize it as JSON
func NewUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	delegateInterceptor := otel.UnaryClientInterceptor(
		global.TraceProvider().Tracer("ai.traceable"),
	)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// GRPC interceptors do not support request/response parsing so the only way to
		// achieve it is by wrapping the invoker (where we can still access the current
		// span).
		wrappedInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			span := trace.SpanFromContext(ctx)
			reqBody, err := internal.MarshalMessageableJSON(req)
			if len(reqBody) > 0 && err == nil {
				span.SetAttribute("grpc.request.body", string(reqBody))
			}

			err = invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return err
			}

			resBody, err := internal.MarshalMessageableJSON(reply)
			if len(resBody) > 0 && err == nil {
				span.SetAttribute("grpc.response.body", string(resBody))
			}

			return err
		}

		return delegateInterceptor(ctx, method, req, reply, cc, wrappedInvoker, opts...)
	}
}
