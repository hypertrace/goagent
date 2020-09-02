// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/traceableai/goagent/examples/grpc/helloworld"
	"github.com/traceableai/goagent/examples/internal"
	traceablegrpc "github.com/traceableai/goagent/otel/grpc"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
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
	internal.InitTracer("grpc-server")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			traceablegrpc.WrapUnaryServerInterceptor(
				otelgrpc.UnaryServerInterceptor(global.TraceProvider().Tracer("ai.traceable")),
			),
		),
	)

	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
