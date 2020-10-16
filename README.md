[![codecov](https://codecov.io/gh/Traceableai/goagent/branch/master/graph/badge.svg?token=MM5BVNGPKE)](https://codecov.io/gh/Traceableai/goagent)

# Go Agent

`goagent` provides a set of complementary instrumentation features for collecting relevant data to be processed by Hypertrace.

## Getting started

`goagent` does not require any particular setup for its instrumentations but it does need to be declared along with standard instrumentation.

- [Getting started with opencensus](instrumentation/opencensus/README.md#getting-started)
- [Getting started with opentelemetry](instrumentation/opentelemetry/README.md#getting-started)

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
