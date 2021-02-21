package hypertrace

import (
	"context"
	"testing"

	"github.com/hypertrace/goagent/config"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

type carrier struct {
	m map[string]string
}

func (c carrier) Get(string) string { return "" }

func (c carrier) Set(key string, value string) {
	c.m[key] = value
}

func TestPropagationFormats(t *testing.T) {
	cfg := config.Load()
	cfg.PropagationFormats = []config.PropagationFormat{
		config.PropagationFormat_B3,
		config.PropagationFormat_TRACECONTEXT,
	}
	Init(cfg)
	tracer := otel.Tracer("b3")
	ctx, _ := tracer.Start(context.Background(), "test")
	propagator := otel.GetTextMapPropagator()
	c := carrier{make(map[string]string)}
	propagator.Inject(ctx, c)
	_, ok := c.m["x-b3-traceid"]
	assert.True(t, ok)
	_, ok = c.m["traceparent"]
	assert.True(t, ok)
}
