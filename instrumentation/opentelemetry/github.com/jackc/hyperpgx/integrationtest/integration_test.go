//go:build integration
// +build integration

package integrationtest

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/jackc/hyperpgx"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal/tracetesting"
	"github.com/stretchr/testify/assert"

	apitrace "go.opentelemetry.io/otel/trace"
)

func TestQuerySuccess(t *testing.T) {
	conn, err := hyperpgx.Connect(context.Background(), "postgres://root:123456@localhost:5432")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	_, flusher := tracetesting.InitTracer()

	select {
	// We timeout after 5 secs of trying ping the DB.
	case <-time.After(time.Duration(5) * time.Second):
		t.Fatal("Unable to ping the DB")
	default:
		err := conn.Ping(context.Background())
		if err != nil {
			t.Logf("Failed to ping the DB: %v", err)
		}
		time.Sleep(time.Duration(200) * time.Millisecond)
	}

	var n int
	err = conn.QueryRow(context.Background(), "SELECT 1 WHERE 1 = $1", 1).Scan(&n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	assert.Equal(t, 1, n)

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, "db:query", span.Name())
	assert.Equal(t, apitrace.SpanKindClient, span.SpanKind())

	attrs := tracetesting.LookupAttributes(span.Attributes())
	assert.Equal(t, "SELECT 1 WHERE 1 = $1", attrs.Get("db.statement").AsString())
	assert.Equal(t, "postgres", attrs.Get("db.system").AsString())
	assert.Equal(t, "root", attrs.Get("db.user").AsString())
	assert.Equal(t, "localhost", attrs.Get("net.peer.name").AsString())
	assert.Equal(t, "5432", attrs.Get("net.peer.port").AsString())
	assert.False(t, attrs.Has("error"))
}
