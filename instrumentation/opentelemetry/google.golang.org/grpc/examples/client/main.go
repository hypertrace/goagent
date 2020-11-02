// +build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	traceablegrpc "github.com/hypertrace/goagent/instrumentation/opentelemetry/google.golang.org/grpc"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/google.golang.org/grpc/examples"
	pb "github.com/hypertrace/goagent/instrumentation/opentelemetry/google.golang.org/grpc/examples/helloworld"
	"go.opencensus.io/trace"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	closer := examples.InitTracer("grpc-client")
	defer closer()

	ctx, span := trace.StartSpan(
		context.Background(),
		"client-bootstrap",
		trace.WithSampler(trace.AlwaysSample()),
	)
	defer span.End()

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			traceablegrpc.WrapUnaryClientInterceptor(
				otelgrpc.UnaryClientInterceptor(global.TraceProvider().Tracer("ai.traceable")),
			),
		),
	)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %v", r.GetMessage())
}
