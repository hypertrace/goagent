package filter

import (
	"testing"

	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter/result"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestMultiFilterEmpty(t *testing.T) {
	f := NewMultiFilter()
	res := f.Evaluate(nil)
	assert.False(t, res.Block)
}

func TestMultiFilterStopsAfterTrue(t *testing.T) {
	tCases := map[string]struct {
		expectedURLAndHeadersFilterResult bool
		expectedBodyFilterResult          bool
		expectedFilterResult              bool
		multiFilter                       *MultiFilter
	}{
		"Evaluate multi filter": {
			expectedFilterResult: true,
			multiFilter: NewMultiFilter(
				mock.Filter{
					Evaluator: func(span sdk.Span) result.FilterResult {
						return result.FilterResult{}
					},
				},
				mock.Filter{
					Evaluator: func(span sdk.Span) result.FilterResult {
						return result.FilterResult{Block: true, ResponseStatusCode: 403}
					},
				},
				mock.Filter{
					Evaluator: func(span sdk.Span) result.FilterResult {
						assert.Fail(t, "should not be called")
						return result.FilterResult{}
					},
				},
			),
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			res := tCase.multiFilter.Evaluate(nil)
			assert.Equal(t, tCase.expectedFilterResult, res.Block)
		})
	}
}
