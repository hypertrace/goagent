# Go Agent

![test](https://github.com/hypertrace/goagent/workflows/test/badge.svg)
[![codecov](https://codecov.io/gh/hypertrace/goagent/branch/master/graph/badge.svg)](https://codecov.io/gh/hypertrace/goagent)

`goagent` provides a set of complementary instrumentation features for collecting relevant data to be processed by [Hypertrace](https://hypertrace.org).

## Getting started

Setting up Go Agent can be done with a few lines:

```go
import "github.com/hypertrace/goagent/config"

...

func main() {
    cfg := config.Load()
    cfg.ServiceName = config.String("myservice")

    shutdown := hypertrace.Init(cfg)
    defer shutdown
}
```

Config values can be declared in config file, env variables or code. For further information about config check [this section](config/README.md).

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

#### Options

##### Filter
[Filtering](sdk/filter/README.md) can be added as part of options. Multiple filters can be added and they will be run in sequence until a filter returns true (request is blocked), or all filters are run.

```go

// ...

    r.Handle("/foo/{bar}", hyperhttp.NewHandler(
        fooHandler,
        "/foo/{bar}",
        hyperhttp.WithFilter(filter.NewMultiFilter(filter1, filter2)),
    ))

// ...

````

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
        http.DefaultTransport,
    ),
}

req, _ := http.NewRequest("GET", "http://example.com", nil)

res, err := client.Do(req)

// ...
```

### Running HTTP examples

In terminal 1 run the client:

```bash
go run ./examples/http-client/main.go
```

In terminal 2 run the server:

```bash
go run ./examples/http-server/main.go
```


## Gin-Gonic Server

Gin server instrumentation relies on adding the `hypergin.Middleware` middleware to the gin server. 
```go
r := gin.Default()

cfg := config.Load()
cfg.ServiceName = config.String("http-gin-server")

flusher := hypertrace.Init(cfg)
defer flusher()

r.Use(hypergin.Middleware())
```

To run an example gin server with the hypertrace middleware: 
```bash
go run ./examples/gin-server/main.go
```

Then make a request to `localhost:8080/ping`

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

#### Options

##### Filter
[Filtering](sdk/filter/README.md) can be added as part of options. Multiple filters can be added and they will be run in sequence until a filter returns true (request is blocked), or all filters are run.

```go

// ...

    grpc.UnaryInterceptor(
        hypergrpc.UnaryServerInterceptor(
            hypergrpc.WithFilter(filter.NewMultiFilter(filter1, filter2))
        ),
    ),

// ...

````

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
go run ./examples/grpc-client/main.go
```

In terminal 2 run the server:

```bash
go run ./examples/grpc-server/main.go
```

## Other instrumentations

- [database/hypersql](instrumentation/hypertrace/database/hypersql)
- [github.com/gorilla/hypermux](instrumentation/hypertrace/github.com/gorilla/hypermux)

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

### Releasing

Run `./release.sh <version_number>` (`<version_number>` should follow semver, e.g. `1.2.3`). The script will change the hardcoded version, commit it, push a tag and prepare the hardcoded version for the next release. After that go to the releases page and draft a new release based on the new tag.

### Further Reference

Read more about `goagent` in the 'Yet Another [Go Agent](https://blog.hypertrace.org/blog/yet-another-go-agent/)' blog post. 
