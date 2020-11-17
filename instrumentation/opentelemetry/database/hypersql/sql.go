package hypersql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	"github.com/ngrok/sqlmw"
	"go.opentelemetry.io/otel/api/global"

	"reflect"
)

var regMu sync.Mutex

type interceptor struct {
	sqlmw.NullInterceptor
	defaultAttributes map[string]string
}

const tracerName = "github.com/hypertrace/goagent"

func (in *interceptor) StmtQueryContext(ctx context.Context, conn driver.StmtQueryContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	ctx, span := global.TraceProvider().Tracer(tracerName).Start(ctx, "sql/query")
	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	span.SetAttribute("db.statement", query)
	defer span.End()

	rows, err := conn.QueryContext(ctx, args)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return rows, err
}

func (in *interceptor) StmtExecContext(ctx context.Context, conn driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	ctx, span := global.TraceProvider().Tracer(tracerName).Start(ctx, "sql/exec")
	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	span.SetAttribute("db.statement", query)
	defer span.End()

	rows, err := conn.ExecContext(ctx, args)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return rows, err
}

func (in *interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	ctx, span := global.TraceProvider().Tracer(tracerName).Start(ctx, "sql/query")
	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}
	span.SetAttribute("db.statement", query)
	defer span.End()

	rows, err := conn.QueryContext(ctx, query, args)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return rows, err
}

func (in *interceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	ctx, span := global.TraceProvider().Tracer(tracerName).Start(ctx, "sql/exec")
	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	span.SetAttribute("db.statement", query)
	defer span.End()

	rows, err := conn.ExecContext(ctx, query, args)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return rows, err
}

func (in *interceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (driver.Tx, error) {
	ctx, span := global.TraceProvider().Tracer(tracerName).Start(ctx, "sql/begin_transaction")
	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	defer span.End()

	tx, err := conn.BeginTx(ctx, txOpts)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return tx, err
}

func (in *interceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (driver.Stmt, error) {
	ctx, span := global.TraceProvider().Tracer(tracerName).Start(ctx, "sql/prepare")
	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	defer span.End()

	tx, err := conn.PrepareContext(ctx, query)
	if err != nil {
		span.RecordError(ctx, err)
	}

	return tx, err
}

func (in *interceptor) TxCommit(ctx context.Context, tx driver.Tx) error {
	ctx, span := global.TraceProvider().Tracer(tracerName).Start(ctx, "sql/commit")
	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	defer span.End()

	err := tx.Commit()
	if err != nil {
		span.RecordError(ctx, err)
	}

	return err
}

func (in *interceptor) TxRollback(ctx context.Context, tx driver.Tx) error {
	ctx, span := global.TraceProvider().Tracer(tracerName).Start(ctx, "sql/rollback")
	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	defer span.End()

	err := tx.Rollback()
	if err != nil {
		span.RecordError(ctx, err)
	}

	return err
}

// driverAttributes returns a list of attributes for given driver
// it relies on reflection to obtain information about the driver
// hidden by driver.Driver interface.
// While using reflection represents an overhead, we only using when
// bootstraping the driver and not anymore after that hence the trade-off
// is acceptable. More so when the other alternative is to do typecast
// across different drivers which will also create a runtime dependency or
// rely on the name assigned to driverName which might not be standard.
//
func getDriverName(d driver.Driver) string {
	elem := reflect.TypeOf(d).Elem()
	pkg, name := elem.PkgPath(), elem.Name()
	return pkg + "." + name
}

// dsnReadWrapper is a driver.Driver that allows to read and parse the DSN.
type dsnReadWrapper struct {
	driver.Driver
	driverName          string
	inDefaultAttributes *map[string]string
}

func (w *dsnReadWrapper) Open(dsn string) (driver.Conn, error) {
	*w.inDefaultAttributes = w.parseDSNAttributes(dsn)
	return w.Driver.Open(dsn)
}

// parseDSNAttributes parses the DSN to obtain attributes like user, ip, port and dbName
func (w *dsnReadWrapper) parseDSNAttributes(dsn string) map[string]string {
	attrs := map[string]string{}
	switch w.driverName {
	case "github.com/mattn/go-sqlite3.SQLiteDriver":
		attrs["db.system"] = "sqlite"
	case "github.com/go-sql-driver/mysql.MySQLDriver":
		if parsedAttrs, err := parseDSN(dsn); err == nil {
			attrs = parsedAttrs
		}
		attrs["db.system"] = "mysql"
	}
	return attrs
}

// Wrap takes a SQL driver and wraps it with Hypertrace instrumentation.
func Wrap(d driver.Driver) driver.Driver {
	driverName := getDriverName(d)
	in := &interceptor{}
	return &dsnReadWrapper{Driver: sqlmw.Driver(d, in), driverName: driverName, inDefaultAttributes: &in.defaultAttributes}
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
