package grpc

import (
	"context"
	"encoding/json"
	"log"
	"net"
	reflect "reflect"

	internal "github.com/traceableai/goagent/otel/grpc/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var _ internal.PersonRegistryServer = server{}

type server struct {
	*internal.UnimplementedPersonRegistryServer
}

func (server) Register(_ context.Context, _ *internal.RegisterRequest) (*internal.RegisterReply, error) {
	return &internal.RegisterReply{Id: 1}, nil
}

// createDialer creates a connection to be used as context dialer in GRPC
// communication.
func createDialer(s *grpc.Server) func(context.Context, string) (net.Conn, error) {
	const bufSize = 1024 * 1024

	listener := bufconn.Listen(bufSize)
	conn := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return conn
}

// jsonEqual compares the JSON from two strings.
func jsonEqual(a, b string) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal([]byte(a), &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(b), &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
