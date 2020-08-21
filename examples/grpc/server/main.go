// +build ignore

package main

import (
	"context"
	"log"
	"net"

	"github.com/traceableai/goagent"
	pb "github.com/traceableai/goagent/examples/grpc/message"
	_ "github.com/traceableai/goagent/otel"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedMessengerServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageReply, error) {
	log.Printf("Received: %v", in.Subject)
	return &pb.MessageReply{Ack: true}, nil
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
	pb.RegisterMessengerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
