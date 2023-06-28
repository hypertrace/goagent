//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/opencensus"
	"github.com/hypertrace/goagent/instrumentation/opencensus/google.golang.org/hypergrpc"
	pb "github.com/hypertrace/goagent/instrumentation/opencensus/google.golang.org/hypergrpc/examples/helloworld"
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
	cfg := config.Load()
	cfg.ServiceName = config.String("grpc-server")

	closer := opencensus.Init(cfg)
	defer closer()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.StatsHandler(hypergrpc.WrapServerHandler(&ocgrpc.ServerHandler{})),
	)

	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
