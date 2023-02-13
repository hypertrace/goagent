package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestStatusCode(t *testing.T) {
	assert.Equal(t, codes.Unauthenticated, StatusCode(401))
	assert.Equal(t, codes.PermissionDenied, StatusCode(403))
	assert.Equal(t, codes.NotFound, StatusCode(404))
	assert.Equal(t, codes.Unauthenticated, StatusCode(407))
	assert.Equal(t, codes.DeadlineExceeded, StatusCode(408))
	assert.Equal(t, codes.FailedPrecondition, StatusCode(412))
	assert.Equal(t, codes.ResourceExhausted, StatusCode(413))
	assert.Equal(t, codes.ResourceExhausted, StatusCode(414))
	assert.Equal(t, codes.ResourceExhausted, StatusCode(429))
	assert.Equal(t, codes.ResourceExhausted, StatusCode(431))
	assert.Equal(t, codes.Unknown, StatusCode(400))
	assert.Equal(t, codes.Unknown, StatusCode(500))
}
