package sql // import "github.com/hypertrace/goagent/sdk/instrumentation/database/sql"

import (
	"context"
	stdSQL "database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	"github.com/hypertrace/goagent/sdk"
	"github.com/ngrok/sqlmw"

	"reflect"
)

var regMu sync.Mutex

type interceptor struct {
	sqlmw.NullInterceptor
	startSpan         sdk.StartSpan
	defaultAttributes map[string]string
}

func (in *interceptor) StmtQueryContext(ctx context.Context, conn driver.StmtQueryContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	ctx, span, end := in.startSpan(ctx, "db:query", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer end()

	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}
	span.SetAttribute("db.statement", query)

	rows, err := conn.QueryContext(ctx, args)
	if err != nil {
		span.SetError(err)
	}

	return rows, err
}

func (in *interceptor) StmtExecContext(ctx context.Context, conn driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	ctx, span, end := in.startSpan(ctx, "db:exec", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer end()

	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}
	span.SetAttribute("db.statement", query)

	rows, err := conn.ExecContext(ctx, args)
	if err != nil {
		span.SetError(err)
	}

	return rows, err
}

func (in *interceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	ctx, span, end := in.startSpan(ctx, "db:query", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer end()

	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}
	span.SetAttribute("db.statement", query)

	rows, err := conn.QueryContext(ctx, query, args)
	if err != nil {
		span.SetError(err)
	}

	return rows, err
}

func (in *interceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	ctx, span, end := in.startSpan(ctx, "db:exec", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer end()

	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}
	span.SetAttribute("db.statement", query)

	rows, err := conn.ExecContext(ctx, query, args)
	if err != nil {
		span.SetError(err)
	}

	return rows, err
}

func (in *interceptor) ConnBeginTx(ctx context.Context, conn driver.ConnBeginTx, txOpts driver.TxOptions) (driver.Tx, error) {
	ctx, span, end := in.startSpan(ctx, "db:begin_transaction", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer end()

	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	tx, err := conn.BeginTx(ctx, txOpts)
	if err != nil {
		span.SetError(err)
	}

	return tx, err
}

func (in *interceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (driver.Stmt, error) {
	ctx, span, end := in.startSpan(ctx, "db:prepare", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer end()

	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	tx, err := conn.PrepareContext(ctx, query)
	if err != nil {
		span.SetError(err)
	}

	return tx, err
}

func (in *interceptor) TxCommit(ctx context.Context, tx driver.Tx) error {
	_, span, end := in.startSpan(ctx, "db:commit", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer end()

	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	err := tx.Commit()
	if err != nil {
		span.SetError(err)
	}

	return err
}

func (in *interceptor) TxRollback(ctx context.Context, tx driver.Tx) error {
	_, span, end := in.startSpan(ctx, "db:rollback", &sdk.SpanOptions{Kind: sdk.SpanKindClient})
	defer end()

	for key, value := range in.defaultAttributes {
		span.SetAttribute(key, value)
	}

	err := tx.Rollback()
	if err != nil {
		span.SetError(err)
	}

	return err
}

// driverAttributes returns a list of attributes for given driver
// it relies on reflection to obtain information about the driver
// hidden by driver.Driver interface.
// While using reflection represents an overhead, we only using when
// bootstrapping the driver and not anymore after that hence the trade-off
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
func Wrap(d driver.Driver, startSpan sdk.StartSpan) driver.Driver {
	driverName := getDriverName(d)
	in := &interceptor{startSpan: startSpan}
	return &dsnReadWrapper{Driver: sqlmw.Driver(d, in), driverName: driverName, inDefaultAttributes: &in.defaultAttributes}
}

// Register initializes and registers the hypersql wrapped database driver
// identified by its driverName. On success it
// returns the generated driverName to use when calling hypersql.Open.
func Register(driverName string, startSpan sdk.StartSpan) (string, error) {
	// retrieve the driver implementation we need to wrap with instrumentation
	db, err := stdSQL.Open(driverName, "")
	if err != nil {
		return "", err
	}
	dri := db.Driver()
	if err = db.Close(); err != nil {
		return "", err
	}

	regMu.Lock()
	defer regMu.Unlock()

	hyperDriverName := fmt.Sprintf("hyper-%s-%d", driverName, len(stdSQL.Drivers()))
	stdSQL.Register(hyperDriverName, Wrap(dri, startSpan))
	return hyperDriverName, nil
}
