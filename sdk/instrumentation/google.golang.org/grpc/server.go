package grpc

import (
	"context"
	"net/http"
	"strings"

	config "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"github.com/hypertrace/goagent/sdk/internal/container"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// Options for gRPC instrumentation
type Options struct {
	Filter filter.Filter
}

// WrapUnaryServerInterceptor returns an interceptor that records the request and response message's body
// and serialize it as JSON
func WrapUnaryServerInterceptor(
	delegateInterceptor grpc.UnaryServerInterceptor,
	spanFromContext sdk.SpanFromContext,
	options *Options,
) grpc.UnaryServerInterceptor {
	defaultAttributes := map[string]string{
		"rpc.system": "grpc",
	}
	if containerID, err := container.GetID(); err == nil {
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
			wrapHandler(info.FullMethod, handler, spanFromContext, defaultAttributes, internalconfig.GetConfig().GetDataCapture(), options),
		)
	}
}

func wrapHandler(
	fullMethod string,
	delegateHandler grpc.UnaryHandler,
	spanFromContext sdk.SpanFromContext,
	defaultAttributes map[string]string,
	dataCaptureConfig *config.DataCapture,
	options *Options,
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

		var filter filter.Filter = &filter.NoopFilter{}
		if options != nil && options.Filter != nil {
			filter = options.Filter
		}

		for key, value := range defaultAttributes {
			span.SetAttribute(key, value)
		}

		pieces := strings.Split(fullMethod[1:], "/")
		span.SetAttribute("rpc.service", pieces[0])
		span.SetAttribute("rpc.method", pieces[1])

		span.SetAttribute("rpc.request.metadata.:method", http.MethodPost)

		setSchemeAttributes(ctx, span)

		reqBody, err := marshalMessageableJSON(req)
		if dataCaptureConfig.RpcBody.Request.Value &&
			len(reqBody) > 0 && err == nil {
			setTruncatedBodyAttribute("request", reqBody, int(dataCaptureConfig.BodyMaxSizeBytes.Value), span)

			if filter.EvaluateBody(span, reqBody) {
				return nil, status.Error(codes.PermissionDenied, "Permission Denied")
			}
		}

		if dataCaptureConfig.RpcMetadata.Request.Value {
			setAttributesFromRequestIncomingMetadata(ctx, span)

			if md, ok := metadata.FromIncomingContext(ctx); ok {
				// TODO: decide what should be passed as URL in GRPC
				if filter.EvaluateURLAndHeaders(span, "", md) {
					return nil, status.Error(codes.PermissionDenied, "Permission Denied")
				}
			}
		}

		res, err := delegateHandler(ctx, req)
		if err != nil {
			return res, err
		}

		resBody, err := marshalMessageableJSON(res)
		if dataCaptureConfig.RpcBody.Response.Value &&
			len(resBody) > 0 && err == nil {
			setTruncatedBodyAttribute("response", resBody, int(dataCaptureConfig.BodyMaxSizeBytes.Value), span)
		}

		return res, err
	}
}

var _ stats.Handler = (*handler)(nil)

type handler struct {
	stats.Handler
	spanFromContext   sdk.SpanFromContext
	defaultAttributes map[string]string
	dataCaptureConfig *config.DataCapture
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (s *handler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	defer s.Handler.HandleRPC(ctx, rs)

	span := s.spanFromContext(ctx)
	if span.IsNoop() {
		// isNoop means either the span is not sampled or there was no span
		// in the request context which means this Handler is not used
		// inside an instrumented Handler, hence we just invoke the delegate
		// handler.
		return
	}

	switch rs := rs.(type) {
	case *stats.Begin:
		for key, value := range s.defaultAttributes {
			span.SetAttribute(key, value)
		}
	case *stats.InPayload:
		body, err := marshalMessageableJSON(rs.Payload)
		if len(body) == 0 || err != nil {
			return
		}

		if rs.IsClient() && s.dataCaptureConfig.RpcBody.Response.Value {
			setTruncatedBodyAttribute("response", body, int(s.dataCaptureConfig.BodyMaxSizeBytes.Value), span)
		} else if !rs.IsClient() && s.dataCaptureConfig.RpcBody.Request.Value {
			setTruncatedBodyAttribute("request", body, int(s.dataCaptureConfig.BodyMaxSizeBytes.Value), span)
		}
	case *stats.InHeader:
		if rs.IsClient() && s.dataCaptureConfig.RpcMetadata.Response.Value {
			setAttributesFromMetadata("response", rs.Header, span)
		} else if !rs.IsClient() && s.dataCaptureConfig.RpcMetadata.Request.Value {
			setAttributesFromMetadata("request", rs.Header, span)
		}
	case *stats.InTrailer:
		if rs.IsClient() && s.dataCaptureConfig.RpcMetadata.Response.Value {
			setAttributesFromMetadata("response", rs.Trailer, span)
		} else if !rs.IsClient() && s.dataCaptureConfig.RpcMetadata.Request.Value {
			setAttributesFromMetadata("request", rs.Trailer, span)
		}
	case *stats.OutPayload:
		body, err := marshalMessageableJSON(rs.Payload)
		if len(body) == 0 || err != nil {
			return
		}

		if rs.IsClient() && s.dataCaptureConfig.RpcBody.Request.Value {
			setTruncatedBodyAttribute("request", body, int(s.dataCaptureConfig.BodyMaxSizeBytes.Value), span)
		} else if !rs.IsClient() && s.dataCaptureConfig.RpcBody.Response.Value {
			setTruncatedBodyAttribute("response", body, int(s.dataCaptureConfig.BodyMaxSizeBytes.Value), span)
		}
	case *stats.OutHeader:
		if rs.IsClient() && s.dataCaptureConfig.RpcMetadata.Request.Value {
			setAttributesFromMetadata("request", rs.Header, span)
		} else if !rs.IsClient() && s.dataCaptureConfig.RpcMetadata.Response.Value {
			setAttributesFromMetadata("response", rs.Header, span)
		}
	case *stats.OutTrailer:
		if rs.IsClient() && s.dataCaptureConfig.RpcMetadata.Request.Value {
			setAttributesFromMetadata("request", rs.Trailer, span)
		} else if !rs.IsClient() && s.dataCaptureConfig.RpcMetadata.Response.Value {
			setAttributesFromMetadata("response", rs.Trailer, span)
		}
	}
}

func (s *handler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	ctx = s.Handler.TagRPC(ctx, rti)
	span := s.spanFromContext(ctx)
	if span.IsNoop() {
		// isNoop means either the span is not sampled or there was no span
		// in the request context which means this Handler is not used
		// inside an instrumented Handler, hence we just invoke the delegate
		// handler.
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
	if containerID, err := container.GetID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}

	return &handler{
		Handler:           delegate,
		spanFromContext:   spanFromContext,
		defaultAttributes: defaultAttributes,
		dataCaptureConfig: internalconfig.GetConfig().GetDataCapture(),
	}
}

func setSchemeAttributes(ctx context.Context, span sdk.Span) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return
	}
	scheme := "http"
	// https://github.com/grpc/grpc-go/blob/ebfe3be62a82434bc83fd7b36410141a603a96be/peer/peer.go#L35-L36
	if peer.AuthInfo != nil {
		scheme = "https"
	}

	span.SetAttribute("rpc.request.metadata.:scheme", scheme)
}
