package hyperpgx

import (
	"context"
	"database/sql/driver"

	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	"github.com/hypertrace/goagent/sdk"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
)

var _ Conn = (*pgx.Conn)(nil)

type Conn interface {
	pgxtype.Querier
	driver.Pinger

	// QueryFunc executes sql with args. For each row returned by the query the values will scanned into the elements of
	// scans and f will be called. If any row fails to scan or f returns an error the query will be aborted and the error
	// will be returned.
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

var _ Conn = (*wrappedConn)(nil)

type wrappedConn struct {
	delegate *pgx.Conn
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
	ctx, span, closer := opentelemetry.StartSpan(ctx, "exec", &sdk.SpanOptions{Kind: sdk.Client})
	defer closer()

	span.SetAttribute("db.statement", query)

	rows, err := w.delegate.Query(ctx, query, optionsAndArgs...)
	if err != nil {
		span.SetError(err)
	}

	return rows, err
}

func (w *wrappedConn) QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row {
	ctx, span, closer := opentelemetry.StartSpan(ctx, "exec", &sdk.SpanOptions{Kind: sdk.Client})
	defer closer()

	span.SetAttribute("db.statement", sql)

	return &wrappedRow{delegate: w.delegate.QueryRow(ctx, sql, optionsAndArgs...), span: span}
}

func (w *wrappedConn) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	ctx, span, closer := opentelemetry.StartSpan(ctx, "exec", &sdk.SpanOptions{Kind: sdk.Client})
	defer closer()

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
	ctx, span, closer := opentelemetry.StartSpan(ctx, "exec", &sdk.SpanOptions{Kind: sdk.Client})
	defer closer()

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

var _ Conn = (*wrappedConn)(nil)

func WrapConnection(c *pgx.Conn) Conn {
	return &wrappedConn{c}
}
