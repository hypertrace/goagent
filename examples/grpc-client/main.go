//go:build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/hypertrace/goagent/config"
	pb "github.com/hypertrace/goagent/examples/internal/helloworld"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address     = "localhost:50151"
	defaultName = "world"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = config.String("grpc-client")
	cfg.Reporting.Endpoint = config.String("localhost:5442")
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP

	flusher := hypertrace.Init(cfg)
	defer flusher()

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(hypergrpc.UnaryClientInterceptor()),
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
		log.Printf("could not greet: %v", err)
	} else {
		log.Printf("Greeting: %v", r.GetMessage())
	}
	// some time to flush the spans
	time.Sleep(2000 * time.Millisecond)
}
