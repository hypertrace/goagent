package filter

import (
	"testing"

	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestNoOpFilter(t *testing.T) {
	f := NoOpFilter{}
	assert.False(t, f.EvaluateURL(nil, ""))
	assert.False(t, f.EvaluateHeaders(nil, nil))
	assert.False(t, f.EvaluateBody(nil, nil))
}

func TestMultiFilterEmpty(t *testing.T) {
	f := NewMultiFilter()
	assert.False(t, f.EvaluateURL(nil, ""))
	assert.False(t, f.EvaluateHeaders(nil, nil))
	assert.False(t, f.EvaluateBody(nil, nil))
}

func TestMultiFilterStopsAfterTrue(t *testing.T) {
	tCases := map[string]struct {
		expectedURLFilterResult     bool
		expectedHeadersFilterResult bool
		expectedBodyFilterResult    bool
		multiFilter                 *MultiFilter
	}{
		"URL multi filter": {
			expectedURLFilterResult:     true,
			expectedHeadersFilterResult: false,
			expectedBodyFilterResult:    false,
			multiFilter: NewMultiFilter(
				mock.Filter{
					URLEvaluator: func(span sdk.Span, url string) bool {
						return false
					},
				},
				mock.Filter{
					URLEvaluator: func(span sdk.Span, url string) bool {
						return true
					},
				},
				mock.Filter{
					URLEvaluator: func(span sdk.Span, url string) bool {
						assert.Fail(t, "should not be called")
						return false
					},
				},
			),
		},
		"Headers multi filter": {
			expectedURLFilterResult:     false,
			expectedHeadersFilterResult: true,
			expectedBodyFilterResult:    false,
			multiFilter: NewMultiFilter(
				mock.Filter{
					HeadersEvaluator: func(span sdk.Span, headers map[string][]string) bool {
						return false
					},
				},
				mock.Filter{
					HeadersEvaluator: func(span sdk.Span, headers map[string][]string) bool {
						return true
					},
				},
				mock.Filter{
					HeadersEvaluator: func(span sdk.Span, headers map[string][]string) bool {
						assert.Fail(t, "should not be called")
						return false
					},
				},
			),
		},
		"Body multi filter": {
			expectedURLFilterResult:     false,
			expectedHeadersFilterResult: false,
			expectedBodyFilterResult:    true,
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
	}

	for name, tCase := range tCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tCase.expectedURLFilterResult, tCase.multiFilter.EvaluateURL(nil, ""))
			assert.Equal(t, tCase.expectedHeadersFilterResult, tCase.multiFilter.EvaluateHeaders(nil, nil))
			assert.Equal(t, tCase.expectedBodyFilterResult, tCase.multiFilter.EvaluateBody(nil, nil))
		})
	}
}
