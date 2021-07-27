package hypersql

// highly inspired in https://github.com/openzipkin-contrib/zipkin-go-sql/blob/master/driver_test.go

import (
	"context"
	"database/sql"
	"testing"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry/internal"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	apitrace "go.opentelemetry.io/otel/trace"
)

func createDB(t *testing.T) (*sql.DB, func() []sdktrace.ReadOnlySpan) {
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
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, "db:query", span.Name())
	assert.Equal(t, apitrace.SpanKindClient, span.SpanKind())

	attrs := internal.LookupAttributes(span.Attributes())
	assert.Equal(t, "SELECT 1 WHERE 1 = ?", attrs.Get("db.statement").AsString())
	assert.Equal(t, "sqlite", attrs.Get("db.system").AsString())
	assert.False(t, attrs.Has("error"))

	db.Close()
}

func TestExecSuccess(t *testing.T) {
	ctx := context.Background()

	db, flusher := createDB(t)

	sqlStmt := `
		drop table if exists foo;
		create table foo (id integer not null primary key, name text);
		delete from foo;
	`

	_, err := db.ExecContext(ctx, sqlStmt)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, "db:exec", span.Name())
	assert.Equal(t, apitrace.SpanKindClient, span.SpanKind())

	attrs := internal.LookupAttributes(span.Attributes())
	assert.False(t, attrs.Has("error"))

	db.Close()
}

func TestTxWithCommitSuccess(t *testing.T) {
	ctx := context.Background()

	db, flusher := createDB(t)

	sqlStmt := `
	drop table if exists foo;
	create table foo (id integer not null primary key, name text);
	delete from foo;
`

	_, err := db.ExecContext(ctx, sqlStmt)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec("1", "こんにちわ世界")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	tx.Commit()

	spans := flusher()
	assert.Equal(t, 5, len(spans))

	assert.Equal(t, "db:exec", spans[0].Name())
	assert.Equal(t, "db:begin_transaction", spans[1].Name())
	assert.Equal(t, "db:prepare", spans[2].Name())
	assert.Equal(t, "db:exec", spans[3].Name())
	assert.Equal(t, "db:commit", spans[4].Name())

	for i := 0; i < 5; i++ {
		assert.Equal(t, apitrace.SpanKindClient, spans[i].SpanKind())
		attrs := internal.LookupAttributes(spans[i].Attributes())
		assert.False(t, attrs.Has("error"))
	}

	db.Close()
}

func TestTxWithRollbackSuccess(t *testing.T) {
	ctx := context.Background()

	db, flusher := createDB(t)

	sqlStmt := `
	drop table if exists foo;
	create table foo (id integer not null primary key, name text);
	delete from foo;
`

	_, err := db.ExecContext(ctx, sqlStmt)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec("1", "こんにちわ世界")
	tx.Rollback()

	spans := flusher()
	assert.Equal(t, 5, len(spans))

	assert.Equal(t, "db:exec", spans[0].Name())
	assert.Equal(t, "db:begin_transaction", spans[1].Name())
	assert.Equal(t, "db:prepare", spans[2].Name())
	assert.Equal(t, "db:exec", spans[3].Name())
	assert.Equal(t, "db:rollback", spans[4].Name())

	for i := 0; i < 5; i++ {
		assert.Equal(t, apitrace.SpanKindClient, spans[i].SpanKind())
		attrs := internal.LookupAttributes(spans[i].Attributes())
		assert.False(t, attrs.Has("error"))
	}

	db.Close()
}
