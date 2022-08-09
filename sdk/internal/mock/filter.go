package mock

import (
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter/result"
)

type Filter struct {
	URLAndHeadersEvaluator func(span sdk.Span, url string, headers map[string][]string) result.FilterResult
	BodyEvaluator          func(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult
	Evaluator              func(span sdk.Span, url string, body []byte, headers map[string][]string) result.FilterResult
}

func (f Filter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
	if f.URLAndHeadersEvaluator == nil {
		return result.FilterResult{}
	}
	return f.URLAndHeadersEvaluator(span, url, headers)
}

func (f Filter) EvaluateBody(span sdk.Span, body []byte, headers map[string][]string) result.FilterResult {
	if f.BodyEvaluator == nil {
		return result.FilterResult{}
	}
	return f.BodyEvaluator(span, body, headers)
}

func (f Filter) Evaluate(span sdk.Span, url string, body []byte, headers map[string][]string) result.FilterResult {
	if f.Evaluator == nil {
		return result.FilterResult{}
	}
	return f.Evaluator(span, url, body, headers)
}
