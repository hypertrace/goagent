package grpc

import (
	"context"

	"github.com/traceableai/goagent/instrumentation/internal"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc"
)

// WrapUnaryClientInterceptor returns an interceptor that records the request and response message's body
// and serialize it as JSON
func WrapUnaryClientInterceptor(delegateInterceptor grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// GRPC interceptors do not support request/response parsing so the only way to
		// achieve it is by wrapping the invoker (where we can still access the current
		// span).
		wrappedInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			span := trace.SpanFromContext(ctx)
			if _, isNoop := span.(trace.NoopSpan); isNoop {
				// isNoop means either the span is not sampled or there was no span
				// in the request context which means this invoker is not used
				// inside an instrumented invoker, hence we just invoke the delegate
				// round tripper.
				return invoker(ctx, method, req, reply, cc, opts...)
			}

			if containerID, err := internal.GetContainerID(); err == nil {
				span.SetAttribute("container_id", containerID)
			}

			reqBody, err := marshalMessageableJSON(req)
			if len(reqBody) > 0 && err == nil {
				span.SetAttribute("grpc.request.body", string(reqBody))
			}

			setAttributesFromOutgoingMetadata(ctx, span)

			err = invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return err
			}

			resBody, err := marshalMessageableJSON(reply)
			if len(resBody) > 0 && err == nil {
				span.SetAttribute("grpc.response.body", string(resBody))
			}

			return err
		}

		return delegateInterceptor(ctx, method, req, reply, cc, wrappedInvoker, opts...)
	}
}
