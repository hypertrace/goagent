package internal

import (
	"go.opentelemetry.io/otel/label"
)

// LookupAttributes allows to lookup an attribute by key. This
// is useful for testing purposes.
type LookupAttributes []label.KeyValue

// Get returns the value of an attribute by key or returns
// an empty one if it can't find it.
func (a LookupAttributes) Get(key string) label.Value {
	for _, kv := range a {
		if string(kv.Key) == key {
			return kv.Value
		}
	}

	return label.Value{}
}

// Has returns true if it can find the attribute by key otherwise false
func (a LookupAttributes) Has(key string) bool {
	for _, kv := range a {
		if string(kv.Key) == key {
			return true
		}
	}

	return false
}
