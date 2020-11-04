# Rationale

## Why wrapping OTel interceptor and not using a chain of interceptors?

GRPC client interceptor is simple, you can think of it as an extension for the OTel interceptor (if OTel interceptor allowed extension). It wraps the OTel interceptor to be able to access its span (in the context) and enrich the trace with the request/reply body and metadata. There are two options to achieve it:

```go
// wrapping the OTel interceptor
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
```

or

```go
// pass the traceable interceptor along with the OTel's one
conn, err := grpc.Dial(
	address,
	grpc.WithInsecure(),
	grpc.WithBlock(),
	grpc.WithChainUnaryInterceptor(
		otelgrpc.UnaryClientInterceptor(myTracer),
            hypergrpc.UnaryClientInterceptor()
	),
)
```

however, according to documentation:

> [...] The first interceptor will be the outer most, while the last interceptor will be the inner most wrapper around the real call.

Unfortunately in our case we need the Traceable's and OTel's interceptors to happen in an specific order, hence we prefferred the composition where the ordering is explicit vs the chain of interceptors where we delegate the responsibility to the user and yet the ordering is not explicit. That said, user can still do 

```go
conn, err := grpc.Dial(
	address,
	grpc.WithInsecure(),
	grpc.WithBlock(),
	grpc.WithChainUnaryInterceptor(
		hypergrpc.WrapUnaryClientInterceptor(
			otelgrpc.UnaryClientInterceptor(myTracer),
        ),
        // other interceptors
	),
)
```
