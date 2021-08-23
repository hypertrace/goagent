package filter

import "github.com/hypertrace/goagent/sdk"

// Filter evaluates whether request should be blocked, `true` blocks the request and `false` continues it.
type Filter interface {
	// Start the filter to allow evaluating requests
	Start()

	// EvaluateURLAndHeaders can be used to evaluate both URL and Headers
	EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) bool

	// EvaluateBody can be used to evaluate the body content
	EvaluateBody(span sdk.Span, body []byte) bool
}
