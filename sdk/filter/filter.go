package filter

import "github.com/hypertrace/goagent/sdk"

// Filter runs on request evaluates whether request should be
// blocked
type Filter struct {
	URLEvaluators     []func(span sdk.Span, URL string) bool
	HeadersEvaluators []func(span sdk.Span, headers map[string][]string) bool
	BodyEvaluators    []func(span sdk.Span, body []byte) bool
}

// EvaluateURL runs URL evaluators in the filter
func (f Filter) EvaluateURL(span sdk.Span, url string) bool {
	for _, urlEvaluator := range f.URLEvaluators {
		if urlEvaluator(span, url) {
			return true
		}
	}

	return false
}

// EvaluateHeaders runs headers evaluators in the filter
func (f Filter) EvaluateHeaders(span sdk.Span, headers map[string][]string) bool {
	for _, headersEvaluator := range f.HeadersEvaluators {
		if headersEvaluator(span, headers) {
			return true
		}
	}

	return false
}

// EvaluateBody runs body evaluators in the filter
func (f Filter) EvaluateBody(span sdk.Span, body []byte) bool {
	for _, bodyEvaluator := range f.BodyEvaluators {
		if bodyEvaluator(span, body) {
			return true
		}
	}

	return false
}
