# Go Agent for OpenCensus

Go Agent provides a set of complementary features for OpenCensus instrumentation

## Package net/http

### HTTP server

The server instrumentation relies on the `http.Handler` component of the server declarations.

```go
import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/hypertrace/goagent/instrumentation/opencensus/net/hyperhttp"
	ochttp "go.opencensus.io/plugin/ochttp"
)

func main() {
    // ...

    r := mux.NewRouter()
    r.Handle("/foo/{bar}", &ochttp.Handler{
        Handler: hyperhttp.WrapHandler(http.HandlerFunc(fooHandler)),
    })

    // ...
}
```

### HTTP client

The client instrumentation relies on the `http.Transport` component of the HTTP client in Go.

```go
import (
    "net/http"
    "github.com/hypertrace/goagent/instrumentation/net/hyperhttp"
    "go.opencensus.io/plugin/ochttp"
)

// ...

client := http.Client{
    Transport: &ochttp.Transport{
        Base: hyperhttp.WrapTransport(http.DefaultTransport),
    },
}

req, _ := http.NewRequest("GET", "http://example.com", nil)

res, err := client.Do(req)

// ...
```

### Running HTTP examples

In terminal 1 run the client:

```bash
go run ./net/hyperhttp/examples/client/main.go
```

In terminal 2 run the server:

```bash
go run ./net/hyperhttp/examples/server/main.go
```

## Package google.golang.org/grpc

### GRPC server

The server instrumentation relies on the `grpc.UnaryServerInterceptor` component of the server declarations.

```go
import (
    // ...

    "github.com/hypertrace/goagent/instrumentation/opencensus/google.golang.org/hypergrpc"
    "go.opencensus.io/plugin/ocgrpc"
    "google.golang.org/grpc"
)


server := grpc.NewServer(
    grpc.UnaryInterceptor(
        grpc.StatsHandler(hypergrpc.WrapServerHandler(&ocgrpc.ServerHandler{})),
    ),
)
```

### GRPC client

The client instrumentation relies on the `http.Transport` component of the HTTP client in Go.

```go
import (
    // ...

    "github.com/hypertrace/goagent/instrumentation/google.golang.org/hypergrpc"
    "go.opencensus.io/plugin/ocgrpc"
    "google.golang.org/grpc"
)

func main() {
    // ...
    conn, err := grpc.Dial(
        address,
        grpc.WithInsecure(),
        grpc.WithBlock(),
        grpc.WithStatsHandler(hypergrpc.WrapClientHandler(&ocgrpc.ClientHandler{})),
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
