package mock

import (
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter/result"
)

type Filter struct {
	Evaluator func(span sdk.Span) result.FilterResult
}

func (f Filter) Evaluate(span sdk.Span) result.FilterResult {
	if f.Evaluator == nil {
		return result.FilterResult{}
	}
	return f.Evaluator(span)
}
