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

// Evaluate runs body evaluators for each filter until one returns true
func (m *MultiFilter) Evaluate(span sdk.Span) result.FilterResult {
	for _, f := range (*m).filters {
		filterResult := f.Evaluate(span)
		if filterResult.Block {
			return filterResult
		}
	}
	return result.FilterResult{}
}
