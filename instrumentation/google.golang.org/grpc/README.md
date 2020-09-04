# goagent for google.golang.org/grpc

## Getting started

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

    traceablegrpc "github.com/traceableai/goagent/instrumentation/google.golang.org/grpc"
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

## Running example

In terminal 1 run

```bash
make run-grpc-server-example
```

In terminal 2 run

```bash
make run-grpc-client-example
```
