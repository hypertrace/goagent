package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/database/hypersql"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

const mysqlLoopCount int = 5

// Run docker run mysql before invoking this
func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("http-server-with-mysql")

	flusher := hypertrace.Init(cfg)
	defer flusher()

	// Set up mysql connection
	db := dbConn()
	defer db.Close()

	r := mux.NewRouter()
	fooHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fooHandler(db, w, r)
	}
	r.Handle("/foo", hyperhttp.NewHandler(
		http.HandlerFunc(fooHandlerFunc),
		"/foo",
	))
	// G114 (CWE-676): Use of net/http serve function that has no support for setting timeouts (Confidence: HIGH, Severity: MEDIUM)
	// #nosec G114
	log.Fatal(http.ListenAndServe(":8081", r))
}

type person struct {
	Name string `json:"name"`
}

func fooHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	sBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	p := &person{}
	err = json.Unmarshal(sBody, p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mysqlLoop(db, r.Context())

	<-time.After(20 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", p.Name)))
}

func dbConn() (db *sql.DB) {
	// Explicitly wrap the MySQLDriver driver with hypersql.
	driver := hypersql.Wrap(&mysql.MySQLDriver{})

	// Register our hypersql wrapper as a database driver.
	sql.Register("ht-mysql", driver)

	db, err := sql.Open("ht-mysql", "root:root@tcp(localhost)/")
	if err != nil {
		log.Fatalf("failed to connect the DB: %v", err)
	}

	log.Println("Connecting to db")
	if err != nil {
		log.Println(err.Error())
	}

	_, err2 := db.ExecContext(context.Background(), "CREATE DATABASE IF NOT EXISTS user")
	if err2 != nil {
		log.Println(err2.Error())
	}

	_, err3 := db.ExecContext(context.Background(), "USE user")
	if err3 != nil {
		log.Println(err3.Error())
	}

	_, err4 := db.ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS Persons (FirstName varchar(255), LastName varchar(255), phoneNumber varchar(255), country varchar(255), serviceCount int)")
	if err4 != nil {
		log.Println(err4.Error())
	}

	return db
}

// mysqlLoop takes in a context which may contain tracing headers which are used to correlate to the caller.
// The caller can be an http server handler function for example.
func mysqlLoop(db *sql.DB, ctx context.Context) {
	for i := 0; i < mysqlLoopCount; i++ {
		_, err1 := db.ExecContext(ctx, "INSERT INTO Persons (FirstName, LastName, phoneNumber, country, serviceCount) VALUES ('Alice', 'Tom B. Erichsen', '1234', 'India', 0)")
		if err1 != nil {
			log.Println(err1.Error())
		}
		_, err2 := db.ExecContext(ctx, "SELECT * FROM Persons")
		if err2 != nil {
			log.Println(err2.Error())
		}

		_, err3 := db.ExecContext(ctx, "DELETE FROM Persons WHERE FirstName = 'Alice'")
		if err3 != nil {
			log.Println(err3.Error())
		}
	}
}
