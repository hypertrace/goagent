package grpc

import (
	"context"

	otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// NewUnaryServerInterceptor returns an interceptor that records the request and response message's body
// and serialize it as JSON
func NewUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	delegateInterceptor := otel.UnaryServerInterceptor(
		global.TraceProvider().Tracer("ai.traceable"),
	)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// GRPC interceptors do not support request/response parsing so the only way to
		// achieve it is by wrapping the handler (where we can still access the current
		// span).
		wrappedHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			span := trace.SpanFromContext(ctx)
			reqBody, err := marshalMessageableJSON(req)
			if len(reqBody) > 0 && err == nil {
				span.SetAttribute("grpc.request.body", string(reqBody))
			}

			res, err := handler(ctx, req)
			if err != nil {
				return res, err
			}

			resBody, err := marshalMessageableJSON(res)
			if len(resBody) > 0 && err == nil {
				span.SetAttribute("grpc.response.body", string(resBody))
			}

			return res, err
		}

		return delegateInterceptor(ctx, req, info, wrappedHandler)
	}
}

// We need a marshaller that does not omit the zero (e.g. 0 or false) to not to lose pontetially
// intesting information.
var marshaler = protojson.MarshalOptions{EmitUnpopulated: true}

// marshalMessageableJSON marshals a value that an be cast as proto.Message into JSON.
func marshalMessageableJSON(messageable interface{}) ([]byte, error) {
	if m, ok := messageable.(proto.Message); ok {
		return marshaler.Marshal(m)
	}

	return nil, nil
}
