// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	traceablegrpc "github.com/traceableai/goagent/instrumentation/opencensus/google.golang.org/grpc"
	pb "github.com/traceableai/goagent/instrumentation/opencensus/google.golang.org/grpc/examples/helloworld"
	"github.com/traceableai/goagent/instrumentation/opencensus/internal"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %q", in.Name)
	return &pb.HelloReply{Message: fmt.Sprintf("hello %s", in.Name)}, nil
}

func main() {
	internal.InitTracer()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.UnaryInterceptor(traceablegrpc.UnaryServerInterceptor()),
	)

	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
