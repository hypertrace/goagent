[![codecov](https://codecov.io/gh/Traceableai/goagent/branch/master/graph/badge.svg?token=MM5BVNGPKE)](https://codecov.io/gh/Traceableai/goagent)

# goagent

`goagent` provides a set of complementary instrumentation features for collecting relevant data to be processed by Hypertrace.

## Getting started

`goagent` does not require any particular setup for OpenTelemetry but it does need to be declared along with OpenTelemetry standard instrumentation, mostly relying on the span being created by OpenTelemetry instrumentation.

- [Getting started with net/http](instrumentation/net/http/README.md#getting-started)
- [Getting started with golang.google.org/grpc](instrumentation/google.golang.org/grpc/README.md#getting-started)

## Contributing

### Running tests

Tests can be run by

```bash
make test
```

for unit tests only

```bash
make test-unit
```
