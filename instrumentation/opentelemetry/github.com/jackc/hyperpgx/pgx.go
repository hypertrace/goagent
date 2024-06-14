package hyperpgx // import "github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/jackc/hyperpgx"

import (
	"context"
	"database/sql/driver"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	"github.com/hypertrace/goagent/sdk"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype/pgxtype"
	pgx "github.com/jackc/pgx/v4"
)

var _ PGXConn = (*pgx.Conn)(nil)

// PGXConn contains all public methods included by *pgx.Conn as an attempt to make the instrumentation transparent
// for the user.
type PGXConn interface {
	pgxtype.Querier
	driver.Pinger

	// QueryFunc executes sql with args. For each row returned by the query the values will scanned into the elements of
	// scans and f will be called. If any row fails to scan or f returns an error the query will be aborted and the error
	// will be returned.
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)

	// SendBatch sends all queued queries to the server at once. All queries are run in an implicit transaction unless
	// explicit transaction control statements are executed. The returned BatchResults must be closed before the connection
	// is used again.
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

	// Close closes a connection. It is safe to call Close on a already closed
	// connection.
	Close(ctx context.Context) error
}

var _ PGXConn = (*wrappedConn)(nil)

type wrappedConn struct {
	delegate  *pgx.Conn
	connAttrs map[string]string
}

var _ pgx.Row = (*wrappedRow)(nil)

type wrappedRow struct {
	delegate pgx.Row
	span     sdk.Span
}

func (r *wrappedRow) Scan(dest ...interface{}) error {
	err := r.delegate.Scan(dest...)
	if err != nil {
		r.span.SetError(err)
	}

	return err
}

func (w *wrappedConn) Query(ctx context.Context, query string, optionsAndArgs ...interface{}) (pgx.Rows, error) {
	ctx, span, closer := opentelemetry.StartSpan(ctx, "db:query", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer closer()

	for k, v := range w.connAttrs {
		span.SetAttribute(k, v)
	}
	span.SetAttribute("db.statement", query)

	rows, err := w.delegate.Query(ctx, query, optionsAndArgs...)
	if err != nil {
		span.SetError(err)
	}

	return rows, err
}

func (w *wrappedConn) QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row {
	ctx, span, closer := opentelemetry.StartSpan(ctx, "db:query", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer closer()

	for k, v := range w.connAttrs {
		span.SetAttribute(k, v)
	}
	span.SetAttribute("db.statement", sql)

	return &wrappedRow{delegate: w.delegate.QueryRow(ctx, sql, optionsAndArgs...), span: span}
}

func (w *wrappedConn) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	ctx, span, closer := opentelemetry.StartSpan(ctx, "exec", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer closer()

	for k, v := range w.connAttrs {
		span.SetAttribute(k, v)
	}
	span.SetAttribute("db.statement", sql)

	res, err := w.delegate.Exec(ctx, sql, arguments...)
	if err != nil {
		span.SetError(err)
	}

	return res, err
}

func (w *wrappedConn) Ping(ctx context.Context) error {
	return w.delegate.Ping(ctx)
}

func (w *wrappedConn) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	ctx, span, closer := opentelemetry.StartSpan(ctx, "exec", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer closer()

	for k, v := range w.connAttrs {
		span.SetAttribute(k, v)
	}
	span.SetAttribute("db.statement", sql)

	res, err := w.delegate.QueryFunc(ctx, sql, args, scans, f)
	if err != nil {
		span.SetError(err)
	}

	return res, err
}

func (w *wrappedConn) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return w.delegate.SendBatch(ctx, b)
}

func (w *wrappedConn) Close(ctx context.Context) error {
	return w.delegate.Close(ctx)
}

var _ PGXConn = (*wrappedConn)(nil)

func Connect(ctx context.Context, connString string) (PGXConn, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return conn, err
	}

	connAttrs, err := parseDSN(connString)
	if err == nil {
		connAttrs["db.system"] = "postgres"
	}

	return &wrappedConn{conn, connAttrs}, nil
}
