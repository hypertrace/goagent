// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hypertrace/goagent/config"
	"github.com/hypertrace/goagent/instrumentation/hypertrace"
	"github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc"
	pb "github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc/examples/helloworld"
	sdkgrpc "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"
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

	flusher := hypertrace.Init(cfg)
	defer flusher()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(hypergrpc.UnaryServerInterceptor(&sdkgrpc.Options{})),
	)

	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
