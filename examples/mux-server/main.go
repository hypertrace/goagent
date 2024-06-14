package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/gorilla/hypermux"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("http-mux-server")

	flusher := hypertrace.Init(cfg)
	defer flusher()

	r := mux.NewRouter()
	r.Use(hypermux.NewMiddleware(&sdkhttp.Options{})) // here we use the mux middleware
	r.HandleFunc("/foo", http.HandlerFunc(fooHandler))
	srv := http.Server{
		Addr:              ":8081",
		Handler:           r,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

type person struct {
	Name string `json:"name"`
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := ioutil.ReadAll(r.Body)
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
	if _, err := w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", p.Name))); err != nil {
		log.Printf("error while writing response body: %v", err)
	}
}
