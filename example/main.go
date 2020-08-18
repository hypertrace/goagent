package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/traceableai/goagent"
	_ "github.com/traceableai/goagent/otel"
)

func main() {
	r := mux.NewRouter()
	r.Handle("/foo", goagent.Instrumentation.HTTPHandler(
		http.HandlerFunc(FooHandler),
	))
	log.Fatal(http.ListenAndServe(":8081", r))
}

func FooHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
