package grpc

import (
	"context"
	"strings"

	"github.com/traceableai/goagent/sdk"
	"github.com/traceableai/goagent/sdk/internal"
	"google.golang.org/grpc"
)

// EnrichUnaryServerInterceptor returns an interceptor that records the request and response message's body
// and serialize it as JSON
func EnrichUnaryServerInterceptor(delegateInterceptor grpc.UnaryServerInterceptor, spanFromContext sdk.SpanFromContext) grpc.UnaryServerInterceptor {
	defaultAttributes := map[string]string{
		"rpc.system": "grpc",
	}
	if containerID, err := internal.GetContainerID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}
	if delegateInterceptor == nil {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return wrapHandler(info.FullMethod, handler, spanFromContext, defaultAttributes)(ctx, req)
		}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// GRPC interceptors do not support request/response parsing so the only way to
		// achieve it is by wrapping the handler (where we can still access the current
		// span).
		return delegateInterceptor(
			ctx,
			req,
			info,
			wrapHandler(info.FullMethod, handler, spanFromContext, defaultAttributes),
		)
	}
}

func wrapHandler(
	fullMethod string,
	delegateHandler grpc.UnaryHandler,
	spanFromContext sdk.SpanFromContext,
	defaultAttributes map[string]string,
) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		span := spanFromContext(ctx)
		if span.IsNoop() {
			// isNoop means either the span is not sampled or there was no span
			// in the request context which means this Handler is not used
			// inside an instrumented Handler, hence we just invoke the delegate
			// round tripper.
			return delegateHandler(ctx, req)
		}
		for key, value := range defaultAttributes {
			span.SetAttribute(key, value)
		}

		pieces := strings.Split(fullMethod[1:], "/")
		span.SetAttribute("rpc.service", pieces[0])
		span.SetAttribute("rpc.method", pieces[1])

		reqBody, err := marshalMessageableJSON(req)
		if len(reqBody) > 0 && err == nil {
			span.SetAttribute("rpc.request.body", string(reqBody))
		}

		setAttributesFromRequestIncomingMetadata(ctx, span)

		res, err := delegateHandler(ctx, req)
		if err != nil {
			return res, err
		}

		resBody, err := marshalMessageableJSON(res)
		if len(resBody) > 0 && err == nil {
			span.SetAttribute("rpc.response.body", string(resBody))
		}

		return res, err
	}
}
