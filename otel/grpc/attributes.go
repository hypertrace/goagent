package grpc

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/metadata"
)

func setAttributes(md metadata.MD, span trace.Span) {
	for key, values := range md {
		if len(values) == 1 {
			span.SetAttribute("grpc.request.metadata."+key, values[0])
			continue
		}

		for index, value := range values {
			span.SetAttribute(
				fmt.Sprintf("grpc.request.metadata.%s[%d]", key, index),
				value,
			)
		}
	}
}

func setAttributesFromOutgoingMetadata(ctx context.Context, span trace.Span) {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		setAttributes(md, span)
	}
}

func setAttributesFromIncomingMetadata(ctx context.Context, span trace.Span) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		setAttributes(md, span)
	}
}
