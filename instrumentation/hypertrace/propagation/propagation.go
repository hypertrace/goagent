package propagation

import (
	"context"

	"go.opentelemetry.io/otel"
)

// TextMapCarrier is the storage medium used by a TextMapPropagator.
type TextMapCarrier interface {
	// Get returns the value associated with the passed key.
	Get(key string) string
	// Set stores the key-value pair.
	Set(key string, value string)
}

// InjectTextMap set cross-cutting concerns from the Context into the TextMap carrier.
func InjectTextMap(ctx context.Context, carrier TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

// ExtractTextMap reads cross-cutting concerns from the TextMap carrier into a Context.
func ExtractTextMap(ctx context.Context, carrier TextMapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}
