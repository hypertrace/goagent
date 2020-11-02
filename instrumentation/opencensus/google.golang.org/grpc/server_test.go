package grpc

import (
	"context"
	"encoding/json"
	"log"
	"net"
	reflect "reflect"
	"testing"

	"github.com/hypertrace/goagent/instrumentation/opencensus/google.golang.org/grpc/examples/helloworld"
	"github.com/hypertrace/goagent/instrumentation/opencensus/internal"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var _ helloworld.GreeterServer = server{}

type server struct {
	reply        *helloworld.HelloReply
	err          error
	replyHeader  metadata.MD
	replyTrailer metadata.MD
	*helloworld.UnimplementedGreeterServer
}

func (s server) SayHello(ctx context.Context, _ *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if s.reply == nil && s.err == nil {
		log.Fatal("missing reply or error in server")
	}

	if err := grpc.SetTrailer(ctx, s.replyTrailer); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to send reply trailer")
	}

	if err := grpc.SendHeader(ctx, s.replyHeader); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to send reply headers")
	}

	return s.reply, s.err
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

func TestServerHelloWorldSuccess(t *testing.T) {
	flusher := internal.InitTracer()

	s := grpc.NewServer(
		grpc.StatsHandler(WrapServerHandler(&ocgrpc.ServerHandler{})),
	)
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{
		reply: &helloworld.HelloReply{Message: "Hi Pupo"},
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("test_key", "test_value"))
	_, err = client.SayHello(ctx, &helloworld.HelloRequest{
		Name: "Pupo",
	})
	if err != nil {
		t.Fatalf("call to Register failed: %v", err)
	}

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	// TODO: Make sure this is consistent with other instrumentations, e.g. OTel
	// records helloworld.Greeter/SayHello
	assert.Equal(t, "helloworld.Greeter.SayHello", span.Name)
	assert.Equal(t, "grpc", span.Attributes["rpc.system"].(string))
	assert.Equal(t, "helloworld.Greeter", span.Attributes["rpc.service"].(string))
	assert.Equal(t, "SayHello", span.Attributes["rpc.method"].(string))
	assert.Equal(t, "test_value", span.Attributes["rpc.request.metadata.test_key"].(string))

	expectedBody := "{\"name\":\"Pupo\"}"
	actualBody := span.Attributes["rpc.request.body"].(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect request body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody = "{\"message\":\"Hi Pupo\"}"
	actualBody = span.Attributes["rpc.response.body"].(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect response body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerHelloWorldFails(t *testing.T) {
	flusher := internal.InitTracer()

	s := grpc.NewServer(
		grpc.StatsHandler(WrapServerHandler(&ocgrpc.ServerHandler{})),
	)
	defer s.Stop()

	expectedError := status.Error(codes.InvalidArgument, "invalid argument")
	helloworld.RegisterGreeterServer(s, &server{
		err: expectedError,
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	_, err = client.SayHello(ctx, &helloworld.HelloRequest{
		Name: "Pupo",
	})
	if err == nil {
		t.Fatalf("expected error: %v", expectedError)
	}

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, int32(codes.InvalidArgument), span.Status.Code)
	assert.Equal(t, "invalid argument", span.Status.Message)
}

func BenchmarkServerRequestResponseBodyMarshaling(b *testing.B) {
	internal.InitTracer()

	s := grpc.NewServer(
		grpc.StatsHandler(WrapServerHandler(&ocgrpc.ServerHandler{})),
	)
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{
		reply: &helloworld.HelloReply{Message: "Hello Pupo"},
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		b.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, err = client.SayHello(ctx, &helloworld.HelloRequest{
			Name: "Pupo",
		})

		if err != nil {
			b.Fatalf("call to Register failed: %v", err)
		}
	}
}

func BenchmarkServerRequestDefaultServerHandler(b *testing.B) {
	s := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	)
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{
		reply: &helloworld.HelloReply{Message: "Hello Pupo"},
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		b.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, err = client.SayHello(ctx, &helloworld.HelloRequest{
			Name: "Pupo",
		})

		if err != nil {
			b.Fatalf("call to Register failed: %v", err)
		}
	}
}
