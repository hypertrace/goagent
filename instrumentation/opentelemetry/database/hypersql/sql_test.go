package hypersql

// highly inspired by https://github.com/openzipkin-contrib/zipkin-go-sql/blob/master/driver_test.go

import (
	"database/sql"
	"testing"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/export/trace"
)

func createDB(t *testing.T) (*sql.DB, func() []*trace.SpanData) {
	_, flusher := internal.InitTracer()

	driverName, err := Register("sqlite3")
	if err != nil {
		t.Fatalf("unable to register driver")
	}

	db, err := sql.Open(driverName, "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}

	return db, flusher
}

func TestQuerySuccess(t *testing.T) {
	db, flusher := createDB(t)

	rows, err := db.Query("SELECT 1 WHERE 1 = ?", 1)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var n int
		if err = rows.Scan(&n); err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
	if err = rows.Err(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	spans := flusher()
	if want, have := 1, len(spans); want != have {
		t.Fatalf("unexpected number of spans, want: %d, have: %d", want, have)
	}

	span := spans[0]
	if want, have := "query", spans[0].Name; want != have {
		t.Fatalf("unexpected span name, want: %s, have: %s", want, have)
	}

	attrs := internal.LookupAttributes(span.Attributes)
	assert.False(t, attrs.Has("error"))

	db.Close()
}
