package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/label"
)

func TestAttributeLookupSuccess(t *testing.T) {
	kvs := []label.KeyValue{
		{Key: "abc", Value: label.StringValue("123")},
	}

	attrs := LookupAttributes(kvs)
	assert.Equal(t, "123", attrs.Get("abc").AsString())
	assert.Equal(t, "", attrs.Get("xyz").AsString())
}

func TestHasAttributeSuccess(t *testing.T) {
	kvs := []label.KeyValue{
		{Key: "abc", Value: label.StringValue("123")},
	}

	attrs := LookupAttributes(kvs)
	assert.True(t, attrs.Has("abc"))
	assert.False(t, attrs.Has("xyz"))
}
