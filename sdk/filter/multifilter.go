package filter // import "github.com/hypertrace/goagent/sdk/filter"

import "github.com/hypertrace/goagent/sdk"

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
