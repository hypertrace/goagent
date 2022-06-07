package mock

import "github.com/hypertrace/goagent/sdk"

type Filter struct {
	URLAndHeadersEvaluator func(span sdk.Span, url string, headers map[string][]string) bool
	BodyEvaluator          func(span sdk.Span, body []byte, headers map[string][]string) bool
	Evaluator              func(span sdk.Span, url string, body []byte, headers map[string][]string) bool
}

func (f Filter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) bool {
	return f.URLAndHeadersEvaluator != nil && f.URLAndHeadersEvaluator(span, url, headers)
}

func (f Filter) EvaluateBody(span sdk.Span, body []byte, headers map[string][]string) bool {
	return f.BodyEvaluator != nil && f.BodyEvaluator(span, body, headers)
}

func (f Filter) Evaluate(span sdk.Span, url string, body []byte, headers map[string][]string) bool {
	return f.Evaluator != nil && f.Evaluator(span, url, body, headers)
}
