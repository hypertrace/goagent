# Go Agent for OpenCensus

Go Agent provides a set of complementary features for OpenCensus instrumentation

## Package net/http

### HTTP server

The server instrumentation relies on the `http.Handler` component of the server declarations.

```go
import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/hypertrace/goagent/instrumentation/opencensus/net/htiohttp"
	ochttp "go.opencensus.io/plugin/ochttp"
)

func main() {
    // ...

    r := mux.NewRouter()
    r.Handle("/foo/{bar}", &ochttp.Handler{
        Handler: htiohttp.WrapHandler(http.HandlerFunc(fooHandler)),
    })

    // ...
}
```

### HTTP client

The client instrumentation relies on the `http.Transport` component of the HTTP client in Go.

```go
import (
    "net/http"
    "github.com/hypertrace/goagent/instrumentation/net/htiohttp"
    "go.opencensus.io/plugin/ochttp"
)

// ...

client := http.Client{
    Transport: &ochttp.Transport{
        Base: htiohttp.WrapTransport(http.DefaultTransport),
    },
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
import (
    // ...

    "github.com/hypertrace/goagent/instrumentation/opencensus/google.golang.org/htiogrpc"
    "go.opencensus.io/plugin/ocgrpc"
    "google.golang.org/grpc"
)


server := grpc.NewServer(
    grpc.UnaryInterceptor(
        grpc.StatsHandler(htiogrpc.WrapServerHandler(&ocgrpc.ServerHandler{})),
    ),
)
```

### GRPC client

The client instrumentation relies on the `http.Transport` component of the HTTP client in Go.

```go
import (
    // ...

    "github.com/hypertrace/goagent/instrumentation/google.golang.org/htiogrpc"
    "go.opencensus.io/plugin/ocgrpc"
    "google.golang.org/grpc"
)

func main() {
    // ...
    conn, err := grpc.Dial(
        address,
        grpc.WithInsecure(),
        grpc.WithBlock(),
        grpc.WithStatsHandler(htiogrpc.WrapClientHandler(&ocgrpc.ClientHandler{})),
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
