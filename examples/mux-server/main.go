package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hypertrace/goagent/config"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"

	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/gorilla/hypermux"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("http-mux-server")
	cfg.Reporting.Endpoint = config.String("localhost:5442")
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP

	flusher := hypertrace.Init(cfg)
	defer flusher()

	r := mux.NewRouter()
	r.Use(hypermux.NewMiddleware(&sdkhttp.Options{})) // here we use the mux middleware
	r.HandleFunc("/foo", http.HandlerFunc(fooHandler))
	log.Fatal(http.ListenAndServe(":8081", r))
}

type person struct {
	Name string `json:"name"`
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
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

	<-time.After(300 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", p.Name)))
}
