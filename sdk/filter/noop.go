package filter // import "github.com/hypertrace/goagent/sdk/filter"

import (
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filterutils"
)

// NoopFilter is a filter that always evaluates to false
type NoopFilter struct{}

var _ Filter = NoopFilter{}

// EvaluateURLAndHeaders that always returns false
func (NoopFilter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) filterutils.FilterResult {
	return filterutils.FilterResult{}
}

// EvaluateBody that always returns false
func (NoopFilter) EvaluateBody(span sdk.Span, body []byte, headers map[string][]string) filterutils.FilterResult {
	return filterutils.FilterResult{}
}

// Evaluate that always returns false
func (NoopFilter) Evaluate(span sdk.Span, url string, body []byte, headers map[string][]string) filterutils.FilterResult {
	return filterutils.FilterResult{}
}
