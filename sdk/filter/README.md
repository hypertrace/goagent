# Go Agent Filter
Filtering capability as specified in [Hypertrace filter specification](https://github.com/hypertrace/specification/blob/main/agent/filtering.md) can be added to instrumentations.

Example
```go
// Filter requests containing "foo" in URL
type FooURLFilter struct {
}

func (FooURLFilter) EvaluateURL(span sdk.Span, url string) bool {
	return strings.Contains(url, "foo")
}

func (FooURLFilter) EvaluateHeaders(span sdk.Span, headers map[string][]string) bool {
	return false
}

func (FooURLFilter) EvaluateBody(span sdk.Span, body []byte) bool {
	return false
}
````