package filter

import (
	"testing"

	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestMultiFilterEmpty(t *testing.T) {
	f := NewMultiFilter()
	res, _ := f.EvaluateURLAndHeaders(nil, "", nil)
	assert.False(t, res)
	res, _ = f.EvaluateBody(nil, nil, nil)
	assert.False(t, res)
	res, _ = f.Evaluate(nil, "", nil, nil)
	assert.False(t, res)
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
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) (bool, int32) {
						return false, 0
					},
				},
				mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) (bool, int32) {
						return true, 403
					},
				},
				mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) (bool, int32) {
						assert.Fail(t, "should not be called")
						return false, 0
					},
				},
			),
		},
		"Body multi filter": {
			expectedBodyFilterResult: true,
			multiFilter: NewMultiFilter(
				mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) (bool, int32) {
						return false, 0
					},
				},
				mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) (bool, int32) {
						return true, 403
					},
				},
				mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte, headers map[string][]string) (bool, int32) {
						assert.Fail(t, "should not be called")
						return false, 0
					},
				},
			),
		},
		"Evaluate multi filter": {
			expectedFilterResult: true,
			multiFilter: NewMultiFilter(
				mock.Filter{
					Evaluator: func(span sdk.Span, url string, body []byte, headers map[string][]string) (bool, int32) {
						return false, 0
					},
				},
				mock.Filter{
					Evaluator: func(span sdk.Span, url string, body []byte, headers map[string][]string) (bool, int32) {
						return true, 403
					},
				},
				mock.Filter{
					Evaluator: func(span sdk.Span, url string, body []byte, headers map[string][]string) (bool, int32) {
						assert.Fail(t, "should not be called")
						return false, 0
					},
				},
			),
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			res, _ := tCase.multiFilter.EvaluateURLAndHeaders(nil, "", nil)
			assert.Equal(t, tCase.expectedURLAndHeadersFilterResult, res)
			res, _ = tCase.multiFilter.EvaluateBody(nil, nil, nil)
			assert.Equal(t, tCase.expectedBodyFilterResult, res)
			res, _ = tCase.multiFilter.Evaluate(nil, "", nil, nil)
			assert.Equal(t, tCase.expectedFilterResult, res)
		})
	}
}
