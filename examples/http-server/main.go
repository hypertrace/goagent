// go:build ignore
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("http-server")
	// Switch to OTLP since default reporting endpoint is ZIPKIN
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP
	cfg.Reporting.Endpoint = config.String("localhost:4317")

	flusher := hypertrace.Init(cfg)
	defer flusher()

	r := mux.NewRouter()
	r.Handle("/foo", hyperhttp.NewHandler(
		http.HandlerFunc(fooHandler),
		"/foo",
	))
	r.Handle("/outgoing", hyperhttp.NewHandler(
		http.HandlerFunc(outgoingCallHandler),
		"/outgoing",
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
	invalidUtf8 := string([]byte{0xff, 0xfe, 0xfd})
	w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s %s\"}", p.Name, invalidUtf8)))
}

func outgoingCallHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		Transport: hyperhttp.NewTransport(
			http.DefaultTransport,
		),
	}
	// "In order to correlate the server call to this client call we need to pass the server request's context."
	req, _ := http.NewRequestWithContext(r.Context(), "GET", "https://httpbin.org/get", nil)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	sb := string(body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, sb)
}
