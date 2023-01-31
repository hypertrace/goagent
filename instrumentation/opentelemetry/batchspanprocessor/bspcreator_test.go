package batchspanprocessor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var batchTimeout = time.Duration(200) * time.Millisecond

func TestCreateBsp(t *testing.T) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	require.NoError(t, err)
	bsp := CreateBatchSpanProcessor(true, exporter,
		sdktrace.WithBatchTimeout(batchTimeout))
	assert.NotNil(t, bsp)

	bsp = CreateBatchSpanProcessor(false, exporter,
		sdktrace.WithBatchTimeout(batchTimeout))
	assert.NotNil(t, bsp)
}
