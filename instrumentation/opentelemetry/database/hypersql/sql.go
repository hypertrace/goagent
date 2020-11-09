package hypersql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	"github.com/ngrok/sqlmw"
	"go.opentelemetry.io/otel/api/global"
)

var regMu sync.Mutex

type interceptor struct {
	sqlmw.NullInterceptor
}

func (in *interceptor) StmtQueryContext(ctx context.Context, conn driver.StmtQueryContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	ctx, span := global.TraceProvider().Tracer("org.hypertrace.goagent").Start(ctx, "query")
	span.SetAttribute("query", query)
	defer span.End()

	rows, err := conn.QueryContext(ctx, args)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return rows, err
}

func (in *interceptor) StmtExecContext(ctx context.Context, conn driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	ctx, span := global.TraceProvider().Tracer("org.hypertrace.goagent").Start(ctx, "exec")
	span.SetAttribute("query", query)
	defer span.End()

	rows, err := conn.ExecContext(ctx, args)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return rows, err
}

func (in *interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	ctx, span := global.TraceProvider().Tracer("org.hypertrace.goagent").Start(ctx, "query")
	span.SetAttribute("query", query)
	defer span.End()

	rows, err := conn.QueryContext(ctx, query, args)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return rows, err
}

func (in *interceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	ctx, span := global.TraceProvider().Tracer("org.hypertrace.goagent").Start(ctx, "exec")
	span.SetAttribute("query", query)
	defer span.End()

	rows, err := conn.ExecContext(ctx, query, args)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return rows, err
}

// Wrap takes a SQL driver and wraps it with Hypertrace instrumentation.
func Wrap(d driver.Driver) driver.Driver {
	return sqlmw.Driver(d, new(interceptor))
}

// Register initializes and registers the hypersql wrapped database driver
// identified by its driverName. On success it
// returns the generated driverName to use when calling sql.Open.
func Register(driverName string) (string, error) {
	// retrieve the driver implementation we need to wrap with instrumentation
	db, err := sql.Open(driverName, "")
	if err != nil {
		return "", err
	}
	dri := db.Driver()
	if err = db.Close(); err != nil {
		return "", err
	}

	regMu.Lock()
	defer regMu.Unlock()

	hyperDriverName := fmt.Sprintf("hyper-%s-%d", driverName, len(sql.Drivers()))
	sql.Register(hyperDriverName, Wrap(dri))
	return hyperDriverName, nil
}
