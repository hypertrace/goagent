package grpc

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/metadata"
)

func setAttributesFromMetadata(md metadata.MD, span trace.Span) {
	for key, values := range md {
		if len(values) == 1 {
			span.SetAttribute("rpc.request.metadata."+key, values[0])
			continue
		}

		for index, value := range values {
			span.SetAttribute(
				fmt.Sprintf("rpc.request.metadata.%s[%d]", key, index),
				value,
			)
		}
	}
}

func setAttributesFromOutgoingMetadata(ctx context.Context, span trace.Span) {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		setAttributesFromMetadata(md, span)
	}
}

func setAttributesFromIncomingMetadata(ctx context.Context, span trace.Span) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		setAttributesFromMetadata(md, span)
	}
}
