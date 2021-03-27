# Go Agent

![test](https://github.com/hypertrace/goagent/workflows/test/badge.svg)
[![codecov](https://codecov.io/gh/hypertrace/goagent/branch/master/graph/badge.svg)](https://codecov.io/gh/hypertrace/goagent)

`goagent` provides a set of complementary instrumentation features for collecting relevant data to be processed by [Hypertrace](https://hypertrace.org). 

## Getting started

Setting up Go Agent can be done with a few lines:

```go
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
    sdkhttp "github.com/hypertrace/goagent/sdk/instrumentation/net/http"
)

func main() {
    // ...

    r := mux.NewRouter()
    r.Handle("/foo/{bar}", hyperhttp.NewHandler(
        fooHandler,
        "/foo/{bar}",
        // See Options section
        &sdkhttp.Options{}
    ))

    // ...
}
```

#### Options

##### Filter
[Filtering](sdk/filter/README.md) can be added as part of Options. Multiple filters can be added and they will be run in sequence until a filter returns true (request is blocked), or all filters are run.

```go

// ...

&sdkhttp.Options {
  Filter: filter.NewMultiFilter(filter1, filter2)
}

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
        hypergrpc.UnaryServerInterceptor(&sdkgrpc.Options{}),
    ),
)
```

#### Options
##### Filter
[Filtering](sdk/filter/README.md) can be added as part of Options. Multiple filters can be added and they will be run in sequence until a filter returns true (request is blocked), or all filters are run.

```go

// ...

&sdkhttp.Options {
  Filter: filter.NewMultiFilter(filter1, filter2)
}

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
go run ./instrumentation/hypertrace/google.golang.org/hypergrpc/examples/client/main.go
```

In terminal 2 run the server:

```bash
go run ./instrumentation/hypertrace/google.golang.org/hypergrpc/examples/server/main.go
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

Run `./release.sh <version_number>` (`<version_number>` should follow semver, e.g. `1.2.3`). The script will change the hardcoded version, commit it, push a tag and prepare the hardcoded version for the next release.

### Further Reference

Read more about `goagent` in the 'Yet Another [Go Agent](https://blog.hypertrace.org/blog/yet-another-go-agent/)' blog post. 
