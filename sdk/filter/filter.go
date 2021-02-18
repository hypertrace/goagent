package filter

// URLFilter filters based on request URL
type URLFilter func(URL string) bool

// HeadersFilter filters based on request headers
type HeadersFilter func(headers map[string][]string) bool

// BodyFilter filters based on request body
type BodyFilter func(body []byte) bool
