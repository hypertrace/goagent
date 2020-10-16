package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traceableai/goagent/sdk/google.golang.org/grpc/internal/helloworld"
	"github.com/traceableai/goagent/sdk/internal/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func makeMockUnaryServerInterceptor(mockSpans *[]*mock.Span) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		span := &mock.Span{}
		ctx = mock.ContextWithSpan(ctx, span)
		*mockSpans = append(*mockSpans, span)
		return handler(ctx, req)
	}
}

func TestServerInterceptorHelloWorldSuccess(t *testing.T) {
	spans := []*mock.Span{}
	mockUnaryInterceptor := makeMockUnaryServerInterceptor(&spans)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			WrapUnaryServerInterceptor(mockUnaryInterceptor, mock.SpanFromContext),
		),
	)
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{})

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

	assert.Equal(t, 1, len(spans))

	span := spans[0]

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

	expectedBody = "{\"message\":\"Hello Pupo\"}"
	actualBody = span.Attributes["rpc.response.body"].(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect response body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerHandlerHelloWorldSuccess(t *testing.T) {
	mockHandler := &mockHandler{}

	s := grpc.NewServer(
		grpc.StatsHandler(WrapStatsHandler(mockHandler, mock.SpanFromContext)),
	)
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{})

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

	assert.Equal(t, 1, len(mockHandler.Spans))

	span := mockHandler.Spans[0]

	assert.Equal(t, "grpc", span.Attributes["rpc.system"].(string))
	assert.Equal(t, "helloworld.Greeter", span.Attributes["rpc.service"].(string))
	assert.Equal(t, "SayHello", span.Attributes["rpc.method"].(string))
	assert.Equal(t, "test_value", span.Attributes["rpc.request.metadata.test_key"].(string))
	assert.Equal(t, "test_value", span.Attributes["rpc.request.metadata.test_key"].(string))

	expectedBody := "{\"name\":\"Pupo\"}"
	actualBody := span.Attributes["rpc.request.body"].(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect request body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody = "{\"message\":\"Hello Pupo\"}"
	actualBody = span.Attributes["rpc.response.body"].(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect response body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}
}
