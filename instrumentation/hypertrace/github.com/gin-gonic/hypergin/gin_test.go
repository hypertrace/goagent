package hypergin

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	m := Middleware()
	assert.NotNil(t, func(*gin.Context) {}, m)
}
