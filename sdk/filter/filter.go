package filter

import "github.com/hypertrace/goagent/sdk"

// Filter evaluates whether request should be blocked, `true` blocks the request and `false` continues it.
type Filter interface {
	// EvaluateURLAndHeaders can be used to evaluate both URL and Headers
	EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) bool

	// EvaluateBody can be used to evaluate the body content
	EvaluateBody(span sdk.Span, body []byte) bool
}

// NoopFilter is a filter that always evaluates to false
type NoopFilter struct {
}

var _ Filter = (*NoopFilter)(nil)

// EvaluateURL that always returns false
func (f *NoopFilter) EvaluateURL(span sdk.Span, url string) bool {
	return false
}

// EvaluateHeaders that always returns false
func (f *NoopFilter) EvaluateHeaders(span sdk.Span, headers map[string][]string) bool {
	return false
}

// EvaluateURLAndHeaders that always returns false
func (f *NoopFilter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) bool {
	return false
}

// EvaluateBody that always returns false
func (f *NoopFilter) EvaluateBody(span sdk.Span, body []byte) bool {
	return false
}

// MultiFilter encapsulates multiple filters
type MultiFilter struct {
	filters []Filter
}

var _ Filter = (*MultiFilter)(nil)

// NewMultiFilter creates a new MultiFilter
func NewMultiFilter(filter ...Filter) *MultiFilter {
	return &MultiFilter{filters: filter}
}

// EvaluateURLAndHeaders runs URL and headers evaluation for each filter until one returns true
func (m *MultiFilter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) bool {
	for _, f := range (*m).filters {
		if f.EvaluateURLAndHeaders(span, url, headers) {
			return true
		}
	}
	return false
}

// EvaluateBody runs body evaluators for each filter until one returns true
func (m *MultiFilter) EvaluateBody(span sdk.Span, body []byte) bool {
	for _, f := range (*m).filters {
		if f.EvaluateBody(span, body) {
			return true
		}
	}
	return false
}
