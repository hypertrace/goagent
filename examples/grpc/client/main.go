// +build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/traceableai/goagent"
	pb "github.com/traceableai/goagent/examples/grpc/message"
	_ "github.com/traceableai/goagent/otel"
	"google.golang.org/grpc"
)

const (
	address        = "localhost:50051"
	defaultSubject = ""
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(goagent.Instrumentation.GRPCInterceptor.UnaryClient()),
		grpc.WithStreamInterceptor(goagent.Instrumentation.GRPCInterceptor.StreamClient()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMessengerClient(conn)

	// Contact the server and print out its response.
	subject := defaultSubject
	if len(os.Args) > 1 {
		subject = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendMessage(ctx, &pb.MessageRequest{Subject: subject})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Ack: %v", r.GetAck())
}
