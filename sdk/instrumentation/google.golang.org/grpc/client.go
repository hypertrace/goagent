package grpc // import "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"

import (
	"context"
	"strings"

	"github.com/hypertrace/goagent/sdk"
	codes "github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter"
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"github.com/hypertrace/goagent/sdk/internal/container"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// WrapUnaryClientInterceptor returns an interceptor that records the request and response message's body
// and serialize it as JSON.
func WrapUnaryClientInterceptor(
	delegateInterceptor grpc.UnaryClientInterceptor,
	spanFromContext sdk.SpanFromContext,
	options *Options,
	spanAttributes map[string]string) grpc.UnaryClientInterceptor {
	var filter filter.Filter = &filter.NoopFilter{}
	if options != nil && options.Filter != nil {
		filter = options.Filter
	}

	defaultAttributes := map[string]string{
		"rpc.system": "grpc",
	}
	for k, v := range spanAttributes {
		defaultAttributes[k] = v
	}
	if containerID, err := container.GetID(); err == nil {
		defaultAttributes["container_id"] = containerID
	}

	dataCaptureConfig := internalconfig.GetConfig().GetDataCapture()

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var header metadata.MD
		var trailer metadata.MD

		// GRPC interceptors do not support request/response parsing so the only way to
		// achieve it is by wrapping the invoker (where we can still access the current
		// span).
		wrappedInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			span := spanFromContext(ctx)
			if span.IsNoop() || span == nil {
				// isNoop means either the span is not sampled or there was no span
				// in the request context which means this invoker is not used
				// inside an instrumented invoker, hence we just invoke the delegate
				// round tripper.
				return invoker(ctx, method, req, reply, cc, opts...)
			}
			for key, value := range defaultAttributes {
				span.SetAttribute(key, value)
			}

			pieces := strings.Split(method[1:], "/")
			span.SetAttribute("rpc.service", pieces[0])
			span.SetAttribute("rpc.method", pieces[1])

			reqBody, err := marshalMessageableJSON(req)
			if dataCaptureConfig.RpcBody.Request.Value && len(reqBody) > 0 && err == nil {
				setTruncatedBodyAttribute("request", reqBody, int(dataCaptureConfig.BodyMaxSizeBytes.Value), span)
			}

			if dataCaptureConfig.RpcMetadata.Request.Value {
				setAttributesFromRequestOutgoingMetadata(ctx, span)
			}

			fr := filter.Evaluate(span)
			if fr.Block {
				statusText := StatusText(int(fr.ResponseStatusCode))
				statusCode := StatusCode(int(fr.ResponseStatusCode))
				span.SetStatus(codes.StatusCodeError, statusText)
				span.SetAttribute("rpc.grpc.status_code", statusCode)
				return status.Error(statusCode, statusText)
			} else if fr.Decorations != nil {
				if md, ok := metadata.FromOutgoingContext(ctx); ok {
					for _, header := range fr.Decorations.RequestHeaderInjections {
						md.Append(header.Key, header.Value)
						span.SetAttribute("rpc.request.metadata."+header.Key, header.Value)
					}
					ctx = metadata.NewIncomingContext(ctx, md)
				}
			}

			err = invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return err
			}

			if dataCaptureConfig.RpcMetadata.Response.Value {
				setAttributesFromMetadata("response", header, span)
				setAttributesFromMetadata("response", trailer, span)
			}

			resBody, err := marshalMessageableJSON(reply)
			if dataCaptureConfig.RpcBody.Response.Value && len(resBody) > 0 && err == nil {
				setTruncatedBodyAttribute("response", resBody, int(dataCaptureConfig.BodyMaxSizeBytes.Value), span)
			}

			return err
		}

		// Even if user pases a header or trailer the data is being populated
		// in all the headers and trailers registered.
		opts = append(opts, grpc.Header(&header), grpc.Trailer(&trailer))

		return delegateInterceptor(ctx, method, req, reply, cc, wrappedInvoker, opts...)
	}
}
