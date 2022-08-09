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
	res := f.EvaluateURLAndHeaders(nil, "", nil)
	assert.False(t, res.Block)
	res = f.EvaluateBody(nil, nil, nil)
	assert.False(t, res.Block)
	res = f.Evaluate(nil, "", nil, nil)
	assert.False(t, res.Block)
}

func TestMultiFilterStopsAfterTrue(t *testing.T) {
	tCases := map[string]struct {
		expectedURLAndHeadersFilterResult bool
		expectedBodyFilterResult          bool
		expectedFilterResult              bool
		multiFilter                       *MultiFilter
	}{
		"URL and Headers multi filter": {
			expectedURLAndHeadersFilterResult: true,
			expectedBodyFilterResult:          false,
			expectedFilterResult:              false,
			multiFilter: NewMultiFilter(
				mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
						return result.FilterResult{}
					},
				},
				mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
						return result.FilterResult{Block: true, ResponseStatusCode: 403}
					},
				},
				mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
						assert.Fail(t, "should not be called")
						return result.FilterResult{}
					},
				},
			),
		},
		"Body multi filter": {
			expectedBodyFilterResult: true,
			multiFilter: NewMultiFilter(
				mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult {
						return result.FilterResult{}
					},
				},
				mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult {
						return result.FilterResult{Block: true, ResponseStatusCode: 403}
					},
				},
				mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult {
						assert.Fail(t, "should not be called")
						return result.FilterResult{}
					},
				},
			),
		},
		"Evaluate multi filter": {
			expectedFilterResult: true,
			multiFilter: NewMultiFilter(
				mock.Filter{
					Evaluator: func(span sdk.Span, url string, body []byte, headers map[string][]string) result.FilterResult {
						return result.FilterResult{}
					},
				},
				mock.Filter{
					Evaluator: func(span sdk.Span, url string, body []byte, headers map[string][]string) result.FilterResult {
						return result.FilterResult{Block: true, ResponseStatusCode: 403}
					},
				},
				mock.Filter{
					Evaluator: func(span sdk.Span, url string, body []byte, headers map[string][]string) result.FilterResult {
						assert.Fail(t, "should not be called")
						return result.FilterResult{}
					},
				},
			),
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			res := tCase.multiFilter.EvaluateURLAndHeaders(nil, "", nil)
			assert.Equal(t, tCase.expectedURLAndHeadersFilterResult, res.Block)
			res = tCase.multiFilter.EvaluateBody(nil, nil, nil)
			assert.Equal(t, tCase.expectedBodyFilterResult, res.Block)
			res = tCase.multiFilter.Evaluate(nil, "", nil, nil)
			assert.Equal(t, tCase.expectedFilterResult, res.Block)
		})
	}
}
