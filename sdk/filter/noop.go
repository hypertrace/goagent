package filter // import "github.com/hypertrace/goagent/sdk/filter"

import (
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter/result"
)

// NoopFilter is a filter that always evaluates to false
type NoopFilter struct{}

var _ Filter = NoopFilter{}

// Evaluate that always returns false
func (NoopFilter) Evaluate(span sdk.Span) result.FilterResult {
	return result.FilterResult{}
}
