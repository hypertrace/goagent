// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

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
		grpc.UnaryInterceptor(serverInterceptor),
		//		grpc.UnaryInterceptor(goagent.Instrumentation.GRPCInterceptor.UnaryServer()),
		//		grpc.StreamInterceptor(goagent.Instrumentation.GRPCInterceptor.StreamServer()),
	)
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Authorization unary interceptor function to handle authorize per RPC call
func serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	t := reflect.TypeOf(req)
	v := reflect.ValueOf(req)
	k := t.Kind()
	fmt.Println("Type ", t)
	fmt.Println("Value ", v)
	fmt.Println("Kind ", k)

	if reflect.ValueOf(req).Kind() == reflect.Ptr {
		r := reflect.ValueOf(req).Elem()
		v := reflect.ValueOf(r)
		fmt.Println("Number of fields", v.NumField())
		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("Field:%d type:%T value:%v\n", i, v.Field(i), v.Field(i))
		}
	}

	// Calls the handler
	h, err := handler(ctx, req)

	// Logging with grpclog (grpclog.LoggerV2)
	log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)

	return h, err
}
