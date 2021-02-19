package filter

import "github.com/hypertrace/goagent/sdk"

// Filter evaluates whether request should be blocked
type Filter interface {
	EvaluateURL(span sdk.Span, url string) bool
	EvaluateHeaders(span sdk.Span, headers map[string][]string) bool
	EvaluateBody(span sdk.Span, body []byte) bool
}

// NoOpFilter is a filter that always evaluates to false
type NoOpFilter struct {
}

// EvaluateURL that always returns false
func (f NoOpFilter) EvaluateURL(span sdk.Span, url string) bool {
	return false
}

// EvaluateHeaders that always returns false
func (f NoOpFilter) EvaluateHeaders(span sdk.Span, headers map[string][]string) bool {
	return false
}

// EvaluateBody that always returns false
func (f NoOpFilter) EvaluateBody(span sdk.Span, body []byte) bool {
	return false
}

// MultiFilter encapsulates multiple filters
type MultiFilter struct {
	filters []Filter
}

// NewMultiFilter creates a new MultiFilter
func NewMultiFilter(filter ...Filter) *MultiFilter {
	return &MultiFilter{filters: filter}
}

// EvaluateURL runs url evaluation for each filter until one returns true
func (m *MultiFilter) EvaluateURL(span sdk.Span, url string) bool {
	for _, f := range (*m).filters {
		if f.EvaluateURL(span, url) {
			return true
		}
	}
	return false
}

// EvaluateHeaders runs headers evaluation for each filter until one returns true
func (m *MultiFilter) EvaluateHeaders(span sdk.Span, headers map[string][]string) bool {
	for _, f := range (*m).filters {
		if f.EvaluateHeaders(span, headers) {
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
