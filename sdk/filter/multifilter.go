package filter // import "github.com/hypertrace/goagent/sdk/filter"

import (
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter/result"
)

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
func (m *MultiFilter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
	for _, f := range (*m).filters {
		filterResult := f.EvaluateURLAndHeaders(span, url, headers)
		if filterResult.Block {
			return filterResult
		}
	}
	return result.FilterResult{}
}

// EvaluateBody runs body evaluators for each filter until one returns true
func (m *MultiFilter) EvaluateBody(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult {
	for _, f := range (*m).filters {
		filterResult := f.EvaluateBody(span, body, headers)
		if filterResult.Block {
			return filterResult
		}
	}
	return result.FilterResult{}
}

// Evaluate runs body evaluators for each filter until one returns true
func (m *MultiFilter) Evaluate(span sdk.Span, url string, body []byte, headers map[string][]string) result.FilterResult {
	for _, f := range (*m).filters {
		filterResult := f.Evaluate(span, url, body, headers)
		if filterResult.Block {
			return filterResult
		}
	}
	return result.FilterResult{}
}
