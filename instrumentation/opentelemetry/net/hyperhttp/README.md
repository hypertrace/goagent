# goagent for net/http

## Getting started

### HTTP server

The server instrumentation relies on the `http.Handler` component of the server declarations.

```go
import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/hypertrace/goagent/instrumentation/net/hyperhttp"
    otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
)

func main() {
    // ...

    r := mux.NewRouter()
    r.Handle("/foo/{bar}", otelhttp.NewHandler(
        hyperhttp.WrapHandler(fooHandler),
        "/foo/{bar}",
    ))

    // ...
}
```

### HTTP client

The client instrumentation relies on the `http.Transport` component of the HTTP client in Go.

```go
import (
    "net/http"
    "github.com/hypertrace/goagent/instrumentation/net/hyperhttp"
    otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
)

// ...

client := http.Client{
    Transport: otelhttp.NewTransport(
        hyperhttp.WrapDefaultTransport(),
        ...
    ),
}

req, _ := http.NewRequest("GET", "http://example.com", nil)

res, err := client.Do(req)

// ...
```

## Running example

In terminal 1 run

```bash
make run-http-server-example
```

In terminal 2 run

```bash
make run-http-client-example
```
