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
func (m *MultiFilter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) (bool, int32) {
	for _, f := range (*m).filters {
		block, blockingStatusCode := f.EvaluateURLAndHeaders(span, url, headers)
		if block {
			return block, blockingStatusCode
		}
	}
	return false, 0
}

// EvaluateBody runs body evaluators for each filter until one returns true
func (m *MultiFilter) EvaluateBody(span sdk.Span, body []byte, headers map[string][]string) (bool, int32) {
	for _, f := range (*m).filters {
		block, blockingStatusCode := f.EvaluateBody(span, body, headers)
		if block {
			return block, blockingStatusCode
		}
	}
	return false, 0
}

// Evaluate runs body evaluators for each filter until one returns true
func (m *MultiFilter) Evaluate(span sdk.Span, url string, body []byte, headers map[string][]string) (bool, int32) {
	for _, f := range (*m).filters {
		block, blockingStatusCode := f.Evaluate(span, url, body, headers)
		if block {
			return block, blockingStatusCode
		}
	}
	return false, 0
}
