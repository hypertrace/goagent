package hypergrpc

import (
	"context"
	"testing"

	"github.com/hypertrace/goagent/instrumentation/opencensus/google.golang.org/hypergrpc/examples/helloworld"
	"github.com/hypertrace/goagent/instrumentation/opencensus/internal"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestClientHelloWorldSuccess(t *testing.T) {
	flusher := internal.InitTracer()

	s := grpc.NewServer()
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{
		reply:        &helloworld.HelloReply{Message: "Hi Pupo"},
		replyHeader:  metadata.Pairs("test_header_key", "test_header_value"),
		replyTrailer: metadata.Pairs("test_trailer_key", "test_trailer_value"),
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithBlock(),
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
		grpc.WithStatsHandler(WrapClientHandler(&ocgrpc.ClientHandler{})),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("test_key_1", "test_value_1"))
	_, err = client.SayHello(
		ctx,
		&helloworld.HelloRequest{
			Name: "Pupo",
		},
	)
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
	assert.Equal(t, "test_value_1", span.Attributes["rpc.request.metadata.test_key_1"].(string))
	assert.Equal(t, "test_header_value", span.Attributes["rpc.response.metadata.test_header_key"].(string))
	assert.Equal(t, "test_trailer_value", span.Attributes["rpc.response.metadata.test_trailer_key"].(string))

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

func TestClientRegisterPersonFails(t *testing.T) {
	flusher := internal.InitTracer()

	s := grpc.NewServer()
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
		grpc.WithStatsHandler(WrapClientHandler(&ocgrpc.ClientHandler{})),
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

func BenchmarkClientRequestResponseBodyMarshaling(b *testing.B) {
	internal.InitTracer()

	s := grpc.NewServer()
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
		grpc.WithStatsHandler(WrapClientHandler(&ocgrpc.ClientHandler{})),
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

func BenchmarkClientRequestDefaultInterceptor(b *testing.B) {
	internal.InitTracer()

	s := grpc.NewServer()
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
		grpc.WithStatsHandler(WrapClientHandler(&ocgrpc.ClientHandler{})),
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
