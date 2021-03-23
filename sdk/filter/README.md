# Go Agent Filter

Filtering capability as specified in [Hypertrace filter specification](https://github.com/hypertrace/specification/blob/main/agent/filtering.md) can be added to instrumentations.

Example

```go
// Filter requests containing "foo" in URL
type FooURLFilter struct {
}

// Filter evaluates whether request should be blocked, `true` blocks the request and `false` continues it.
func (FooURLFilter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) bool {
	return false
}

// Filter evaluates whether request should be blocked, `true` blocks the request and `false` continues it.
func (FooURLFilter) EvaluateBody(span sdk.Span, body []byte) bool {
	return false
}
```
