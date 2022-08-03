package mock

import "github.com/hypertrace/goagent/sdk"

type Filter struct {
	URLAndHeadersEvaluator func(span sdk.Span, url string, headers map[string][]string) (bool, int32)
	BodyEvaluator          func(span sdk.Span, body []byte, headers map[string][]string) (bool, int32)
	Evaluator              func(span sdk.Span, url string, body []byte, headers map[string][]string) (bool, int32)
}

func (f Filter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) (bool, int32) {
	if f.URLAndHeadersEvaluator == nil {
		return false, 0
	}
	return f.URLAndHeadersEvaluator(span, url, headers)
}

func (f Filter) EvaluateBody(span sdk.Span, body []byte, headers map[string][]string) (bool, int32) {
	if f.BodyEvaluator == nil {
		return false, 0
	}
	return f.BodyEvaluator(span, body, headers)
}

func (f Filter) Evaluate(span sdk.Span, url string, body []byte, headers map[string][]string) (bool, int32) {
	if f.Evaluator == nil {
		return false, 0
	}
	return f.Evaluator(span, url, body, headers)
}
