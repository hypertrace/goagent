# Go Agent Filter

Filtering capability as specified in [Hypertrace filter specification](https://github.com/hypertrace/specification/blob/main/agent/filtering.md) can be added to instrumentations.

Example

```go
// Filter requests containing "foo" in URL
type FooURLFilter struct {
}

// `result.FilterResult` will contain a bool `Block` which when set to `true` means to block the request and `false` to
// continue the request. It also contains `ResponseStatusCode` which is the HTTP status code to return.

// Filter evaluates whether request should be blocked based on the url and headers.
func (FooURLFilter) EvaluateURLAndHeaders(span sdk.Span, url string, headers map[string][]string) result.FilterResult {
	return false
}

// Filter evaluates whether request should be blocked based on the body.
func (FooURLFilter) EvaluateBody(span sdk.Span, body []byte) result.FilterResult {
	return false
}

// Filter evaluates whether request should be blocked based on url, headers and body.
func (FooURLFilter) Evaluate(span sdk.Span, url string, body []byte, headers map[string][]string) result.FilterResult {
	return false
}
```
