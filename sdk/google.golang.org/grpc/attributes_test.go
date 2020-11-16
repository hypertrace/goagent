package grpc

import (
	"testing"

	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestSetScalarAttributeSuccess(t *testing.T) {
	md := metadata.Pairs("key_1", "value_1")
	span := mock.NewSpan()
	setAttributesFromMetadata("request", md, span)

	assert.Equal(t, "value_1", span.ReadAttribute("rpc.request.metadata.key_1").(string))
	assert.Zero(t, span.RemainingAttributes())
}

func TestSetMultivalueAttributeSuccess(t *testing.T) {
	md := metadata.Pairs("key_1", "value_1", "key_1", "value_2")
	span := mock.NewSpan()
	setAttributesFromMetadata("request", md, span)

	assert.Equal(t, "value_1", span.ReadAttribute("rpc.request.metadata.key_1[0]").(string))
	assert.Equal(t, "value_2", span.ReadAttribute("rpc.request.metadata.key_1[1]").(string))
	assert.Zero(t, span.RemainingAttributes())
}
