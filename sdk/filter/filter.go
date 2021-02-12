package filter

import (
	"github.com/hypertrace/goagent/sdk"
)

// Filter is the interface that evaluates
// whether server request should be blocked based
// on request span attributes.
type Filter interface {
	Id() string

	// Evaluate evaluates whether request
	// represented with given request span attributes should be filtered(blocked)
	// The filter may add attributes to provided span
	Evaluate(requestSpanAttributes map[string]string, span sdk.Span) bool
}
