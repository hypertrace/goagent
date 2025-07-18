//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hypertrace/goagent/config"
	pb "github.com/hypertrace/goagent/examples/internal/helloworld"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc"
	"google.golang.org/grpc"
)

const (
	port = ":50151"
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
	cfg := config.Load()
	cfg.ServiceName = config.String("grpc-server")
	cfg.Reporting.Endpoint = config.String("localhost:5442")
	cfg.Reporting.TraceReporterType = config.TraceReporterType_OTLP

	flusher := hypertrace.Init(cfg)
	defer flusher()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(hypergrpc.UnaryServerInterceptor()),
	)

	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
