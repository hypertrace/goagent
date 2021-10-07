package main

// GOOS=linux go build
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("tim-http-server")

	flusher := hypertrace.Init(cfg)
	defer flusher()

	r := mux.NewRouter()
	r.Handle("/foo", hyperhttp.NewHandler(
		http.HandlerFunc(fooHandler),
		"/foo",
	))
	r.Handle("/bigfoo", hyperhttp.NewHandler(
		http.HandlerFunc(repeatedFooHandler),
		"/bigfoo",
	))
	r.Handle("/bigfoorequest", hyperhttp.NewHandler(
		http.HandlerFunc(repeatedFooBigRequestHandler),
		"/bigfoo",
	))
	log.Fatal(http.ListenAndServe(":8081", r))
}

type person struct {
	Name string `json:"name"`
}

type personAndSettings struct {
	Name       string        `json:"name"`
	Iterations int           `json:"iterations"` // 1 iteration is a response body of 25 bytes
	SleepMs    time.Duration `json:"sleep_ms"`
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
		//log.Printf("An error occured: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	<-time.After(300 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"message\": \"Hello %s\"}", p.Name)))
}

func repeatedFooHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	p := &personAndSettings{
		Iterations: 10,
		SleepMs:    30,
	}
	err = json.Unmarshal(sBody, p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	<-time.After(p.SleepMs * time.Millisecond)

	var sb strings.Builder
	sb.WriteString("{")
	for i := 0; i < p.Iterations; i++ {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("\"message_%d\":\"Hello %s\"", i, p.Name))
	}
	sb.WriteString("}")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sb.String()))
}

func repeatedFooBigRequestHandler(w http.ResponseWriter, r *http.Request) {
	sBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	iterations := (len(sBody) / 25) + 1

	<-time.After(5 * time.Millisecond)

	var sb strings.Builder
	sb.WriteString("{")
	for i := 0; i < iterations; i++ {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("\"message_%d\":\"Hello World!\"", i))
	}
	sb.WriteString("}")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sb.String()))
}
