package filter // import "github.com/hypertrace/goagent/sdk/filter"

import "github.com/hypertrace/goagent/sdk"

// NoopFilter is a filter that always evaluates to false
type NoopFilter struct{}

var _ Filter = NoopFilter{}

// EvaluateURLAndHeaders that always returns false
func (NoopFilter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) bool {
	return false
}

// EvaluateBody that always returns false
func (NoopFilter) EvaluateBody(span sdk.Span, body []byte, headers map[string][]string) bool {
	return false
}
