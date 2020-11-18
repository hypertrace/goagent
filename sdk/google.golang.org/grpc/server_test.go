package grpc

import (
	"context"
	"testing"

	"github.com/hypertrace/goagent/sdk/google.golang.org/grpc/internal/helloworld"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func makeMockUnaryServerInterceptor(mockSpans *[]*mock.Span) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		span := mock.NewSpan()
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

	assert.Equal(t, "grpc", span.ReadAttribute("rpc.system").(string))
	assert.Equal(t, "helloworld.Greeter", span.ReadAttribute("rpc.service").(string))
	assert.Equal(t, "SayHello", span.ReadAttribute("rpc.method").(string))
	assert.Equal(t, "test_value", span.ReadAttribute("rpc.request.metadata.test_key").(string))

	expectedBody := "{\"name\":\"Pupo\"}"
	actualBody := span.ReadAttribute("rpc.request.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect request body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody = "{\"message\":\"Hello Pupo\"}"
	actualBody = span.ReadAttribute("rpc.response.body").(string)
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
		grpc.WithUserAgent("test_agent"),
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

	assert.Equal(t, "grpc", span.ReadAttribute("rpc.system").(string))
	assert.Equal(t, "helloworld.Greeter", span.ReadAttribute("rpc.service").(string))
	assert.Equal(t, "SayHello", span.ReadAttribute("rpc.method").(string))
	assert.Equal(t, "test_value", span.ReadAttribute("rpc.request.metadata.test_key").(string))

	assert.Equal(t, "bufnet", span.ReadAttribute("rpc.request.metadata.:authority").(string))
	assert.Equal(t, "application/grpc", span.ReadAttribute("rpc.request.metadata.content-type").(string))
	assert.Contains(t, span.ReadAttribute("rpc.request.metadata.user-agent").(string), "test_agent")

	expectedBody := "{\"name\":\"Pupo\"}"
	actualBody := span.ReadAttribute("rpc.request.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect request body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody = "{\"message\":\"Hello Pupo\"}"
	actualBody = span.ReadAttribute("rpc.response.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect response body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
}
