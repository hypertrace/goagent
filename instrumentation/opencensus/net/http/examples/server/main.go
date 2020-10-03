// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	traceablehttp "github.com/traceableai/goagent/instrumentation/opencensus/net/http"
	"github.com/traceableai/goagent/instrumentation/opencensus/net/http/examples"
	ochttp "go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
)

func main() {
	examples.InitTracer("http-server")

	r := mux.NewRouter()
	r.Handle("/foo", &ochttp.Handler{
		Propagation: &b3.HTTPFormat{},
		Handler:     traceablehttp.EnrichHandler(http.HandlerFunc(fooHandler)),
	})
	log.Fatal(http.ListenAndServe(":8081", r))
}

type person struct {
	Name string `json:"name"`
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p := &person{}
	err = json.Unmarshal(sBody, p)
	if err != nil {
		log.Printf("failed to unmarshal body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", p.Name)))
}
