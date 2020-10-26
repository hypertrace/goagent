package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traceableai/goagent/sdk/internal/mock"
	"google.golang.org/grpc/metadata"
)

func TestSetScalarAttributeSuccess(t *testing.T) {
	md := metadata.Pairs("key_1", "value_1")
	span := mock.NewSpan()
	setAttributesFromMetadata("request", md, span)

	assert.Equal(t, "value_1", span.Attributes["rpc.request.metadata.key_1"].(string))
}

func TestSetMultivalueAttributeSuccess(t *testing.T) {
	md := metadata.Pairs("key_1", "value_1", "key_1", "value_2")
	span := mock.NewSpan()
	setAttributesFromMetadata("request", md, span)

	assert.Equal(t, "value_1", span.Attributes["rpc.request.metadata.key_1[0]"].(string))
	assert.Equal(t, "value_2", span.Attributes["rpc.request.metadata.key_1[1]"].(string))
}
