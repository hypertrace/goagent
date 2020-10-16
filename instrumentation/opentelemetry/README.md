# Go Agent for OpenTelemetry

Go Agent provides a set of complementary features for OpenTelemetry instrumentation

## Package net/http

### HTTP server

The server instrumentation relies on the `http.Handler` component of the server declarations.

```go
import (
    "net/http"

    "github.com/gorilla/mux"
    traceablehttp "github.com/traceableai/goagent/instrumentation/opentelemetry/net/http"
    otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
)

func main() {
    // ...

    r := mux.NewRouter()
    r.Handle("/foo/{bar}", otelhttp.NewHandler(
        traceablehttp.WrapHandler(fooHandler),
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
    traceablehttp "github.com/traceableai/goagent/instrumentation/opentelemetry/net/http"
    otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
)

// ...

client := http.Client{
    Transport: otelhttp.NewTransport(
        traceablehttp.WrapDefaultTransport(),
        ...
    ),
}

req, _ := http.NewRequest("GET", "http://example.com", nil)

res, err := client.Do(req)

// ...
```

### Running HTTP examples

In terminal 1 run the client:

```bash
go run ./net/http/examples/client/main.go
```

In terminal 2 run the server:

```bash
go run ./net/http/examples/server/main.go
```

## Package google.golang.org/grpc

### GRPC server

The server instrumentation relies on the `grpc.UnaryServerInterceptor` component of the server declarations.

```go

server := grpc.NewServer(
    grpc.UnaryInterceptor(
        traceablegrpc.WrapUnaryServerInterceptor(
            otelgrpc.UnaryServerInterceptor(myTracer),
        ),
    ),
)
```

### GRPC client

The client instrumentation relies on the `http.Transport` component of the HTTP client in Go.

```go
import (
    // ...

    traceablegrpc "github.com/traceableai/goagent/instrumentation/opentelemetry/google.golang.org/grpc"
    otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
    "google.golang.org/grpc"
)

func main() {
    // ...
    conn, err := grpc.Dial(
        address,
        grpc.WithInsecure(),
        grpc.WithBlock(),
        grpc.WithUnaryInterceptor(
            traceablegrpc.WrapUnaryClientInterceptor(
                otelgrpc.UnaryClientInterceptor(myTracer),
            ),
        ),
    )
    if err != nil {
        log.Fatalf("could not dial: %v", err)
    }
    defer conn.Close()

    client := pb.NewCustomClient(conn)

    // ...
}
```

### Running GRPC examples

In terminal 1 run the client:

```bash
go run ./google.golang.org/grpc/examples/client/main.go
```

In terminal 2 run the server:

```bash
go run ./google.golang.org/grpc/examples/server/main.go
```
