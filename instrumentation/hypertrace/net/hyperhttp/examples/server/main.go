// +build ignore

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
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
	sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("http-server")

	flusher := hypertrace.Init(cfg)
	defer flusher()

	r := mux.NewRouter()
	r.Handle("/foo", hyperhttp.NewHandler(
		http.HandlerFunc(fooHandler),
		"/foo",
		&sdkhttp.Options{},
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
