package mock

import "github.com/hypertrace/goagent/sdk"

type Filter struct {
	URLEvaluator     func(span sdk.Span, URL string) bool
	HeadersEvaluator func(span sdk.Span, headers map[string][]string) bool
	BodyEvaluator    func(span sdk.Span, body []byte) bool
}

func (f Filter) EvaluateURL(span sdk.Span, url string) bool {
	return f.URLEvaluator != nil && f.URLEvaluator(span, url)
}

func (f Filter) EvaluateHeaders(span sdk.Span, headers map[string][]string) bool {
	return f.HeadersEvaluator != nil && f.HeadersEvaluator(span, headers)
}

func (f Filter) EvaluateBody(span sdk.Span, body []byte) bool {
	return f.BodyEvaluator != nil && f.BodyEvaluator(span, body)
}
