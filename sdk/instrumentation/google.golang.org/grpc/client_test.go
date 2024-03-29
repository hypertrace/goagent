package grpc

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc/internal/helloworld"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func makeMockUnaryClientInterceptor(mockSpans *[]*mock.Span) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		span := mock.NewSpan()
		ctx = mock.ContextWithSpan(ctx, span)
		*mockSpans = append(*mockSpans, span)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func TestUnaryClientHelloWorldSuccess(t *testing.T) {
	spans := []*mock.Span{}

	s := grpc.NewServer()
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{
		replyHeader:  metadata.Pairs("test_header_key", "test_header_value"),
		replyTrailer: metadata.Pairs("test_trailer_key", "test_trailer_value"),
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			WrapUnaryClientInterceptor(
				makeMockUnaryClientInterceptor(&spans),
				mock.SpanFromContext,
				map[string]string{"foo": "bar"},
			),
		),
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
			Name: "Cuchi",
		},
	)
	if err != nil {
		t.Fatalf("call to Register failed: %v", err)
	}

	assert.Equal(t, 1, len(spans))

	span := spans[0]

	assert.Equal(t, "grpc", span.ReadAttribute("rpc.system").(string))
	assert.Equal(t, "helloworld.Greeter", span.ReadAttribute("rpc.service").(string))
	assert.Equal(t, "SayHello", span.ReadAttribute("rpc.method").(string))
	assert.Equal(t, "test_value_1", span.ReadAttribute("rpc.request.metadata.test_key_1").(string))
	assert.Equal(t, "test_header_value", span.ReadAttribute("rpc.response.metadata.test_header_key").(string))
	assert.Equal(t, "test_trailer_value", span.ReadAttribute("rpc.response.metadata.test_trailer_key").(string))
	assert.Equal(t, "application/grpc", span.ReadAttribute("rpc.response.metadata.content-type").(string))
	assert.Equal(t, "bar", span.ReadAttribute("foo").(string))

	expectedBody := "{\"name\":\"Cuchi\"}"
	actualBody := span.ReadAttribute("rpc.request.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect request body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody = "{\"message\":\"Hello Cuchi\"}"
	actualBody = span.ReadAttribute("rpc.response.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect response body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	_ = span.ReadAttribute("container_id") // needed in containarized envs
	assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
}

func TestClientHandlerHelloWorldSuccess(t *testing.T) {
	mockHandler := &mockHandler{}

	s := grpc.NewServer()
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{
		replyHeader:  metadata.Pairs("test_header_key", "test_header_value"),
		replyTrailer: metadata.Pairs("test_trailer_key", "test_trailer_value"),
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(WrapStatsHandler(mockHandler, mock.SpanFromContext)),
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
			Name: "Cuchi",
		},
	)
	if err != nil {
		t.Fatalf("call to Register failed: %v", err)
	}

	assert.Equal(t, 1, len(mockHandler.Spans))

	span := mockHandler.Spans[0]

	assert.Equal(t, "grpc", span.ReadAttribute("rpc.system").(string))
	assert.Equal(t, "helloworld.Greeter", span.ReadAttribute("rpc.service").(string))
	assert.Equal(t, "SayHello", span.ReadAttribute("rpc.method").(string))
	assert.Equal(t, "test_value_1", span.ReadAttribute("rpc.request.metadata.test_key_1").(string))
	assert.Equal(t, "test_header_value", span.ReadAttribute("rpc.response.metadata.test_header_key").(string))
	assert.Equal(t, "test_trailer_value", span.ReadAttribute("rpc.response.metadata.test_trailer_key").(string))

	expectedBody := "{\"name\":\"Cuchi\"}"
	actualBody := span.ReadAttribute("rpc.request.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect request body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody = "{\"message\":\"Hello Cuchi\"}"
	actualBody = span.ReadAttribute("rpc.response.body").(string)
	if ok, err := jsonEqual(expectedBody, actualBody); err == nil {
		assert.True(t, ok, "incorrect response body:\nwant %s,\nhave %s", expectedBody, actualBody)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBodyTruncation(t *testing.T) {
	spans := []*mock.Span{}

	s := grpc.NewServer()
	defer s.Stop()

	helloworld.RegisterGreeterServer(s, &server{
		replyHeader:  metadata.Pairs("test_header_key", "test_header_value"),
		replyTrailer: metadata.Pairs("test_trailer_key", "test_trailer_value"),
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			WrapUnaryClientInterceptor(
				makeMockUnaryClientInterceptor(&spans),
				mock.SpanFromContext,
				map[string]string{"foo": "bar"},
			),
		),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	world40K := strings.Repeat("World", 40000)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("test_key_1", "test_value_1"))
	_, err = client.SayHello(
		ctx,
		&helloworld.HelloRequest{
			Name: world40K,
		},
	)
	if err != nil {
		t.Fatalf("call to Register failed: %v", err)
	}

	assert.Equal(t, 1, len(spans))

	span := spans[0]

	assert.Equal(t, "grpc", span.ReadAttribute("rpc.system").(string))
	assert.Equal(t, "helloworld.Greeter", span.ReadAttribute("rpc.service").(string))
	assert.Equal(t, "SayHello", span.ReadAttribute("rpc.method").(string))
	assert.Equal(t, "test_value_1", span.ReadAttribute("rpc.request.metadata.test_key_1").(string))
	assert.Equal(t, "test_header_value", span.ReadAttribute("rpc.response.metadata.test_header_key").(string))
	assert.Equal(t, "test_trailer_value", span.ReadAttribute("rpc.response.metadata.test_trailer_key").(string))
	assert.Equal(t, "application/grpc", span.ReadAttribute("rpc.response.metadata.content-type").(string))
	assert.Equal(t, "bar", span.ReadAttribute("foo").(string))
	assert.Equal(t, true, span.ReadAttribute("rpc.request.body.truncated").(bool))
	assert.Equal(t, true, span.ReadAttribute("rpc.response.body.truncated").(bool))

	fullBody := fmt.Sprintf("{\"name\":\"%s\"}", world40K)
	expectedBody := fullBody[:131_072]
	actualBody := span.ReadAttribute("rpc.request.body").(string)
	// direct comparison of the body since it will be truncated
	assert.Equal(t, expectedBody, actualBody)

	fullResponseBody := fmt.Sprintf("{\"message\":\"Hello %s\"}", world40K)
	expectedBody = fullResponseBody[:131_072] //"{\"message\":\"Hello Cuchi\"}"
	actualBody = span.ReadAttribute("rpc.response.body").(string)
	// direct comparison of the body since it will be truncated
	assert.Equal(t, expectedBody, actualBody)

	_ = span.ReadAttribute("container_id") // needed in containarized envs
	assert.Zero(t, span.RemainingAttributes(), "unexpected remaining attribute: %v", span.Attributes)
}
