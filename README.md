# goagent

`goagent` provides a set of features that allows the user to instrument applications.

## Running example

### HTTP

Run

```bash
make run-http-server-example
```

once the server is running you can call it using curl:

```bash
curl -i localhost:8081/foo
```

### GRPC

In terminal 1 run

```bash
make run-grpc-server-example
```

In terminal 2 run

```bash
make run-grpc-client-example
```
