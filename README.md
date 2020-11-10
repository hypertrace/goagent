# Go Agent

[![CircleCI](https://circleci.com/gh/hypertrace/goagent/tree/main.svg?style=svg)](https://circleci.com/gh/hypertrace/goagent/tree/main)
[![codecov](https://codecov.io/gh/hypertrace/goagent/branch/master/graph/badge.svg)](https://codecov.io/gh/hypertrace/goagent)

`goagent` provides a set of complementary instrumentation features for collecting relevant data to be processed by [Hypertrace](https://hypertrace.org).

## Getting started

## Package net/hyperhttp

### HTTP server

The server instrumentation relies on the `http.Handler` component of the server declarations.

```go
import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

func main() {
    // ...

    r := mux.NewRouter()
    r.Handle("/foo/{bar}", hyperhttp.NewHandler(
        fooHandler,
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
    "github.com/hypertrace/goagent/instrumentation/hypertrace/net/hyperhttp"
)

// ...

client := http.Client{
    Transport: hyperhttp.NewTransport(
        http.DefaultTransport
    ),
}

req, _ := http.NewRequest("GET", "http://example.com", nil)

res, err := client.Do(req)

// ...
```

### Running HTTP examples

In terminal 1 run the client:

```bash
go run ./instrumentation/hypertrace/net/hyperhttp/examples/client/main.go
```

In terminal 2 run the server:

```bash
go run ./instrumentation/hypertrace/net/hyperhttp/examples/server/main.go
```

## Package google.golang.org/hypergrpc

### GRPC server

The server instrumentation relies on the `grpc.UnaryServerInterceptor` component of the server declarations.

```go

server := grpc.NewServer(
    grpc.UnaryInterceptor(
        hypergrpc.UnaryServerInterceptor(),
    ),
)
```

### GRPC client

The client instrumentation relies on the `http.Transport` component of the HTTP client in Go.

```go
import (
    // ...

    hypergrpc "github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc"
    "google.golang.org/grpc"
)

func main() {
    // ...
    conn, err := grpc.Dial(
        address,
        grpc.WithInsecure(),
        grpc.WithBlock(),
        grpc.WithUnaryInterceptor(
            hypergrpc.UnaryClientInterceptor(),
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
go run ./instrumentation/hypertrace/google.golang.org/grpc/examples/client/main.go
```

In terminal 2 run the server:

```bash
go run ./instrumentation/hypertrace/google.golang.org/grpc/examples/server/main.go
```

## Contributing

### Running tests

Tests can be run with (requires docker)

```bash
make test
```

for unit tests only

```bash
make test-unit
```
