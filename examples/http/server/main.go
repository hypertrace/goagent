// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type person struct {
	Name string `json:"name"`
}

func FooHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p := &person{}
	err = json.Unmarshal(sBody, p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", p.Name)))
}
