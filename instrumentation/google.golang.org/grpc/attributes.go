package grpc

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/metadata"
)

func setAttributesFromMetadata(_type string, md metadata.MD, span trace.Span) {
	for key, values := range md {
		if len(values) == 1 {
			span.SetAttribute(
				fmt.Sprintf("rpc.%s.metadata.%s", _type, key),
				values[0],
			)
			continue
		}

		for index, value := range values {
			span.SetAttribute(
				fmt.Sprintf("rpc.%s.metadata.%s[%d]", _type, key, index),
				value,
			)
		}
	}
}

func setAttributesFromRequestOutgoingMetadata(ctx context.Context, span trace.Span) {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		setAttributesFromMetadata("request", md, span)
	}
}

func setAttributesFromRequestIncomingMetadata(ctx context.Context, span trace.Span) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		setAttributesFromMetadata("request", md, span)
	}
}
