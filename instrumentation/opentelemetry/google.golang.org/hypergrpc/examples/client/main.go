// +build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry"
	"github.com/hypertrace/goagent/instrumentation/opentelemetry/google.golang.org/hypergrpc"
	pb "github.com/hypertrace/goagent/instrumentation/opentelemetry/google.golang.org/hypergrpc/examples/helloworld"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("grpc-client")

	closer := opentelemetry.Init(cfg)
	defer closer()

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			hypergrpc.WrapUnaryClientInterceptor(
				otelgrpc.UnaryClientInterceptor(),
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %v", r.GetMessage())
}
