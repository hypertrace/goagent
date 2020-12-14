# Go Agent for OpenTelemetry

Go Agent provides a set of complementary features for OpenTelemetry instrumentation

## Package net/hyperhttp

### HTTP server

The server instrumentation relies on the `http.Handler` component of the server declarations.

```go
import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/hypertrace/goagent/instrumentation/opentelemetry/net/hyperhttp"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
    "github.com/hypertrace/goagent/instrumentation/opentelemetry/net/hyperhttp"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

### Running HTTP examples

In terminal 1 run the client:

```bash
go run ./net/http/examples/client/main.go
```

In terminal 2 run the server:

```bash
go run ./net/http/examples/server/main.go
```

## Package google.golang.org/hypergrpc

### GRPC server

The server instrumentation relies on the `grpc.UnaryServerInterceptor` component of the server declarations.

```go

server := grpc.NewServer(
    grpc.UnaryInterceptor(
        hypergrpc.WrapUnaryServerInterceptor(
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

    hypergrpc "github.com/hypertrace/goagent/instrumentation/opentelemetry/google.golang.org/hypergrpc"
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
    "google.golang.org/grpc"
)

func main() {
    // ...
    conn, err := grpc.Dial(
        address,
        grpc.WithInsecure(),
        grpc.WithBlock(),
        grpc.WithUnaryInterceptor(
            hypergrpc.WrapUnaryClientInterceptor(
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
go run ./google.golang.org/hypergrpc/examples/client/main.go
```

In terminal 2 run the server:

```bash
go run ./google.golang.org/hypergrpc/examples/server/main.go
```
