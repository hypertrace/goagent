package mock

import (
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filterutils"
)

type Filter struct {
	URLAndHeadersEvaluator func(span sdk.Span, url string, headers map[string][]string) filterutils.FilterResult
	BodyEvaluator          func(span sdk.Span, body []byte, headers map[string][]string) filterutils.FilterResult
	Evaluator              func(span sdk.Span, url string, body []byte, headers map[string][]string) filterutils.FilterResult
}

func (f Filter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) filterutils.FilterResult {
	if f.URLAndHeadersEvaluator == nil {
		return filterutils.FilterResult{}
	}
	return f.URLAndHeadersEvaluator(span, url, headers)
}

func (f Filter) EvaluateBody(span sdk.Span, body []byte, headers map[string][]string) filterutils.FilterResult {
	if f.BodyEvaluator == nil {
		return filterutils.FilterResult{}
	}
	return f.BodyEvaluator(span, body, headers)
}

func (f Filter) Evaluate(span sdk.Span, url string, body []byte, headers map[string][]string) filterutils.FilterResult {
	if f.Evaluator == nil {
		return filterutils.FilterResult{}
	}
	return f.Evaluator(span, url, body, headers)
}
