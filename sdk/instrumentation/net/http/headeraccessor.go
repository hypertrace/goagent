package http // import "github.com/hypertrace/goagent/sdk/instrumentation/net/http"

import "net/http"

// HeaderAccessor allows accessing HTTP header values.
//
// Go packages use varying data structures and conventions for storing header key-values.
// Using this interface allows us our functions to not be tied to a particular such format.
type HeaderAccessor interface {
	Lookup(key string) []string
	ForEachHeader(callback func(key string, values []string) error) error
}

type headerMapAccessor struct {
	header http.Header
}

var _ HeaderAccessor = headerMapAccessor{}

// NewHeaderMapAccessor returns a HeaderAccessor
func NewHeaderMapAccessor(h http.Header) HeaderAccessor {
	return &headerMapAccessor{h}
}

func (a headerMapAccessor) Lookup(key string) []string {
	return a.header[http.CanonicalHeaderKey(key)]
}

func (a headerMapAccessor) ForEachHeader(callback func(key string, values []string) error) error {
	for key, values := range a.header {
		if err := callback(key, values); err != nil {
			return err
		}
	}
	return nil
}
