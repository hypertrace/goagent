package filter // import "github.com/hypertrace/goagent/sdk/filter"

import (
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter/result"
)

// Filter evaluates whether request should be blocked, `true` blocks the request and `false` continues it.
type Filter interface {
	// EvaluateURLAndHeaders can be used to evaluate both URL and headers
	EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) result.FilterResult

	// EvaluateBody can be used to evaluate the body content
	EvaluateBody(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult

	// Evaluate can be used to evaluate URL, headers and body content in one call
	Evaluate(span sdk.Span, url string, body []byte, headers map[string][]string) result.FilterResult
}
