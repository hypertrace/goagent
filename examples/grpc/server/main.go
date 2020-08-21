// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/traceableai/goagent"
	pb "github.com/traceableai/goagent/examples/grpc/helloworld"
	_ "github.com/traceableai/goagent/otel"
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
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(goagent.Instrumentation.GRPCInterceptor.UnaryServer()),
		grpc.StreamInterceptor(goagent.Instrumentation.GRPCInterceptor.StreamServer()),
	)
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}