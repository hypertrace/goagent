package grpc // import "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"

import (
	"context"
	"fmt"

	"github.com/hypertrace/goagent/sdk"
	"google.golang.org/grpc/metadata"
)

func setAttributesFromMetadata(_type string, md metadata.MD, span sdk.Span) {
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

func setAttributesFromRequestOutgoingMetadata(ctx context.Context, span sdk.Span) {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		setAttributesFromMetadata("request", md, span)
	}
}

func setAttributesFromRequestIncomingMetadata(ctx context.Context, span sdk.Span) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		setAttributesFromMetadata("request", md, span)
	}
}
