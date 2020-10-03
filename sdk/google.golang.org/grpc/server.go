package grpc

import (
	"context"

	"github.com/traceableai/goagent/sdk"
	"github.com/traceableai/goagent/sdk/internal"
	"google.golang.org/grpc"
)

// EnrichUnaryServerInterceptor returns an interceptor that records the request and response message's body
// and serialize it as JSON
func EnrichUnaryServerInterceptor(delegateInterceptor grpc.UnaryServerInterceptor, spanFromContext sdk.SpanFromContext) grpc.UnaryServerInterceptor {
	defaultAttributes := make(map[string]string)
	if containerID, err := internal.GetContainerID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// GRPC interceptors do not support request/response parsing so the only way to
		// achieve it is by wrapping the handler (where we can still access the current
		// span).
		wrappedHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			span := spanFromContext(ctx)
			if span.IsNoop() {
				// isNoop means either the span is not sampled or there was no span
				// in the request context which means this Handler is not used
				// inside an instrumented Handler, hence we just invoke the delegate
				// round tripper.
				return handler(ctx, req)
			}
			for key, value := range defaultAttributes {
				span.SetAttribute(key, value)
			}

			reqBody, err := marshalMessageableJSON(req)
			if len(reqBody) > 0 && err == nil {
				span.SetAttribute("rpc.request.body", string(reqBody))
			}

			setAttributesFromRequestIncomingMetadata(ctx, span)

			res, err := handler(ctx, req)
			if err != nil {
				return res, err
			}

			resBody, err := marshalMessageableJSON(res)
			if len(resBody) > 0 && err == nil {
				span.SetAttribute("rpc.response.body", string(resBody))
			}

			return res, err
		}

		return delegateInterceptor(ctx, req, info, wrappedHandler)
	}
}
