// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/traceableai/goagent/instrumentation/opentelemetry/google.golang.org/grpc/examples"
	traceablehttp "github.com/traceableai/goagent/instrumentation/opentelemetry/net/http"
	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
)

func main() {
	examples.InitTracer("http-server")

	r := mux.NewRouter()
	r.Handle("/foo", otelhttp.NewHandler(
		traceablehttp.EnrichHandler(
			http.HandlerFunc(fooHandler),
		),
		"/foo",
	))
	log.Fatal(http.ListenAndServe(":8081", r))
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
