package tracetesting

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestAttributeLookupSuccess(t *testing.T) {
	kvs := []attribute.KeyValue{
		{Key: "abc", Value: attribute.StringValue("123")},
	}

	attrs := LookupAttributes(kvs)
	assert.Equal(t, "123", attrs.Get("abc").AsString())
	assert.Equal(t, "", attrs.Get("xyz").AsString())
}

func TestHasAttributeSuccess(t *testing.T) {
	kvs := []attribute.KeyValue{
		{Key: "abc", Value: attribute.StringValue("123")},
	}

	attrs := LookupAttributes(kvs)
	assert.True(t, attrs.Has("abc"))
	assert.False(t, attrs.Has("xyz"))
}
