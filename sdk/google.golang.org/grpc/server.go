package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/traceableai/goagent/sdk"
	"github.com/traceableai/goagent/sdk/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

// WrapUnaryServerInterceptor returns an interceptor that records the request and response message's body
// and serialize it as JSON
func WrapUnaryServerInterceptor(
	delegateInterceptor grpc.UnaryServerInterceptor,
	spanFromContext sdk.SpanFromContext,
) grpc.UnaryServerInterceptor {
	defaultAttributes := map[string]string{
		"rpc.system": "grpc",
	}
	if containerID, err := internal.GetContainerID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
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

var _ stats.Handler = (*handler)(nil)

type handler struct {
	stats.Handler
	spanFromContext   sdk.SpanFromContext
	defaultAttributes map[string]string
}

func resolveInType(isClient bool) string {
	if isClient {
		return "response"
	}

	return "request"
}

func resolveOutType(isClient bool) string {
	return resolveInType(!isClient)
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (s *handler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	defer s.Handler.HandleRPC(ctx, rs)

	span := s.spanFromContext(ctx)
	if span.IsNoop() {
		// isNoop means either the span is not sampled or there was no span
		// in the request context which means this Handler is not used
		// inside an instrumented Handler, hence we just invoke the delegate
		// round tripper.
		return
	}

	switch rs := rs.(type) {
	case *stats.Begin:
		for key, value := range s.defaultAttributes {
			span.SetAttribute(key, value)
		}
	case *stats.InPayload:
		body, err := marshalMessageableJSON(rs.Payload)
		if len(body) > 0 && err == nil {
			span.SetAttribute(
				fmt.Sprintf("rpc.%s.body", resolveInType(rs.IsClient())),
				string(body),
			)
		}
	case *stats.InHeader:
		setAttributesFromMetadata(resolveInType(rs.IsClient()), rs.Header, span)
	case *stats.InTrailer:
		setAttributesFromMetadata(resolveInType(rs.IsClient()), rs.Trailer, span)
	case *stats.OutPayload:
		body, err := marshalMessageableJSON(rs.Payload)
		if len(body) > 0 && err == nil {
			span.SetAttribute(
				fmt.Sprintf("rpc.%s.body", resolveOutType(rs.IsClient())),
				string(body),
			)
		}
	case *stats.OutHeader:
		setAttributesFromMetadata(resolveOutType(rs.IsClient()), rs.Header, span)
	case *stats.OutTrailer:
		setAttributesFromMetadata(resolveOutType(rs.IsClient()), rs.Trailer, span)
	}
}

func (s *handler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	ctx = s.Handler.TagRPC(ctx, rti)
	span := s.spanFromContext(ctx)
	if span.IsNoop() {
		// isNoop means either the span is not sampled or there was no span
		// in the request context which means this Handler is not used
		// inside an instrumented Handler, hence we just invoke the delegate
		// round tripper.
		return ctx
	}

	pieces := strings.Split(rti.FullMethodName[1:], "/")
	span.SetAttribute("rpc.service", pieces[0])
	span.SetAttribute("rpc.method", pieces[1])

	return ctx
}

// WrapStatsHandler wraps an instrumented StatsHandler and returns a new one that records
// the request/response body and metadata.
func WrapStatsHandler(delegate stats.Handler, spanFromContext sdk.SpanFromContext) stats.Handler {
	defaultAttributes := map[string]string{
		"rpc.system": "grpc",
	}
	if containerID, err := internal.GetContainerID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}

	return &handler{Handler: delegate, spanFromContext: spanFromContext, defaultAttributes: defaultAttributes}
}
