package filter

import (
	"testing"

	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestMultiFilterEmpty(t *testing.T) {
	f := NewMultiFilter()
	assert.False(t, f.EvaluateURLAndHeaders(nil, "", nil))
	assert.False(t, f.EvaluateBody(nil, nil))
	assert.False(t, f.Evaluate(nil))
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
			multiFilter: NewMultiFilter(
				mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) bool {
						return false
					},
				},
				mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) bool {
						return true
					},
				},
				mock.Filter{
					URLAndHeadersEvaluator: func(span sdk.Span, url string, headers map[string][]string) bool {
						assert.Fail(t, "should not be called")
						return false
					},
				},
			),
		},
		"Body multi filter": {
			expectedBodyFilterResult: true,
			multiFilter: NewMultiFilter(
				mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte) bool {
						return false
					},
				},
				mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte) bool {
						return true
					},
				},
				mock.Filter{
					BodyEvaluator: func(span sdk.Span, body []byte) bool {
						assert.Fail(t, "should not be called")
						return false
					},
				},
			),
		},
		"Evaluate multi filter": {
			expectedFilterResult: true,
			multiFilter: NewMultiFilter(
				mock.Filter{
					Evaluator: func(span sdk.ReadbackSpan) bool {
						return false
					},
				},
				mock.Filter{
					Evaluator: func(span sdk.ReadbackSpan) bool {
						return true
					},
				},
				mock.Filter{
					Evaluator: func(span sdk.ReadbackSpan) bool {
						assert.Fail(t, "should not be called")
						return false
					},
				},
			),
		},
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tCase.expectedURLAndHeadersFilterResult, tCase.multiFilter.EvaluateURLAndHeaders(nil, "", nil))
			assert.Equal(t, tCase.expectedBodyFilterResult, tCase.multiFilter.EvaluateBody(nil, nil))
			assert.Equal(t, tCase.expectedFilterResult, tCase.multiFilter.Evaluate(nil))
		})
	}
}
