// +build ignore

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	traceablehttp "github.com/traceableai/goagent/instrumentation/opencensus/net/http"
	"github.com/traceableai/goagent/instrumentation/opencensus/net/http/examples"
	ochttp "go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

type message struct {
	Content string `json:"message"`
}

func main() {
	closer := examples.InitTracer("http-client")
	defer closer()

	ctx, span := trace.StartSpan(
		context.Background(),
		"client-bootstrap",
		trace.WithSampler(trace.AlwaysSample()),
	)
	defer span.End()

	client := http.Client{
		Transport: &ochttp.Transport{
			Base: traceablehttp.WrapTransport(http.DefaultTransport),
		},
	}

	req, err := http.NewRequest("GET", "http://localhost:8081/foo", bytes.NewBufferString(`{"name":"Dave"}`))
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalf("failed to create the request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to perform the request: %v", err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read the response body: %v", err)
	}
	defer res.Body.Close()
	m := &message{}
	err = json.Unmarshal(resBody, m)
	if err != nil {
		log.Fatalf("failed to unmarshal the response body: %v", err)
		return
	}

	fmt.Println(m.Content)
}
