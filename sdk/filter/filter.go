package filter // import "github.com/hypertrace/goagent/sdk/filter"

import (
	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/filter/result"
)

// Filter evaluates whether request should be blocked, `true` blocks the request and `false` continues it.
type Filter interface {
	// Evaluate can be used to evaluate URL, headers and body content in one call
	Evaluate(span sdk.Span) result.FilterResult
}
