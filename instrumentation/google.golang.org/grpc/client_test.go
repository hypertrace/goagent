package grpc

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	grpcinternal "github.com/traceableai/goagent/instrumentation/google.golang.org/grpc/internal"
	"github.com/traceableai/goagent/instrumentation/internal"
	otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestClientRegisterPersonSuccess(t *testing.T) {
	_, flusher := internal.InitTracer()

	s := grpc.NewServer()
	defer s.Stop()

	grpcinternal.RegisterPersonRegistryServer(s, &server{
		reply:        &grpcinternal.RegisterReply{Id: 1},
		replyHeader:  metadata.Pairs("test_header_key", "test_header_value"),
		replyTrailer: metadata.Pairs("test_trailer_key", "test_trailer_value"),
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			WrapUnaryClientInterceptor(
				otel.UnaryClientInterceptor(global.TraceProvider().Tracer("ai.traceable")),
			),
		),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcinternal.NewPersonRegistryClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("test_key_1", "test_value_1"))
	_, err = client.Register(
		ctx,
		&grpcinternal.RegisterRequest{
			Firstname: "Bugs",
			Lastname:  "Bunny",
			Birthdate: &timestamp.Timestamp{Seconds: 1},
			Confirmed: false,
		},
	)
	if err != nil {
		t.Fatalf("call to Register failed: %v", err)
	}

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, "helloworld.PersonRegistry/Register", span.Name)

	attrs := internal.LookupAttributes(span.Attributes)
	assert.Equal(t, "grpc", attrs.Get("rpc.system").AsString())
	assert.Equal(t, "helloworld.PersonRegistry", attrs.Get("rpc.service").AsString())
	assert.Equal(t, "Register", attrs.Get("rpc.method").AsString())
	assert.Equal(t, "test_value_1", attrs.Get("rpc.request.metadata.test_key_1").AsString())
	assert.Equal(t, "test_header_value", attrs.Get("rpc.response.metadata.test_header_key").AsString())
	assert.Equal(t, "test_trailer_value", attrs.Get("rpc.response.metadata.test_trailer_key").AsString())

	expectedBody := "{\"firstname\":\"Bugs\",\"lastname\":\"Bunny\",\"birthdate\":\"1970-01-01T00:00:01Z\",\"confirmed\":false}"
	if ok, err := jsonEqual(expectedBody, attrs.Get("rpc.request.body").AsString()); err == nil {
		assert.True(t, ok)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBody = "{\"id\":\"1\"}"
	if ok, err := jsonEqual(expectedBody, attrs.Get("rpc.response.body").AsString()); err == nil {
		assert.True(t, ok)
	} else {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientRegisterPersonFails(t *testing.T) {
	_, flusher := internal.InitTracer()

	s := grpc.NewServer()
	defer s.Stop()

	expectedError := status.Error(codes.InvalidArgument, "invalid argument")
	grpcinternal.RegisterPersonRegistryServer(s, &server{
		err: expectedError,
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			WrapUnaryClientInterceptor(
				otel.UnaryClientInterceptor(global.TraceProvider().Tracer("ai.traceable")),
			),
		),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcinternal.NewPersonRegistryClient(conn)

	_, err = client.Register(ctx, &grpcinternal.RegisterRequest{
		Firstname: "Bugs",
		Lastname:  "Bunny",
		Birthdate: &timestamp.Timestamp{Seconds: 1},
		Confirmed: false,
	})
	if err == nil {
		t.Fatalf("expected error: %v", expectedError)
	}

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, codes.InvalidArgument, span.StatusCode)
	assert.Equal(t, "invalid argument", span.StatusMessage)
}

func BenchmarkClientRequestResponseBodyMarshaling(b *testing.B) {
	internal.InitTracer()

	s := grpc.NewServer()
	defer s.Stop()

	grpcinternal.RegisterPersonRegistryServer(s, &server{
		reply: &grpcinternal.RegisterReply{Id: 1},
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			WrapUnaryClientInterceptor(
				otel.UnaryClientInterceptor(global.TraceProvider().Tracer("ai.traceable")),
			),
		),
	)
	if err != nil {
		b.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcinternal.NewPersonRegistryClient(conn)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, err = client.Register(ctx, &grpcinternal.RegisterRequest{
			Firstname: "Bugs",
			Lastname:  "Bunny",
			Birthdate: &timestamp.Timestamp{Seconds: int64(n)},
			Confirmed: false,
		})

		if err != nil {
			b.Fatalf("call to Register failed: %v", err)
		}
	}
}

func BenchmarkClientRequestDefaultInterceptor(b *testing.B) {
	tracer, _ := internal.InitTracer()

	s := grpc.NewServer()
	defer s.Stop()

	grpcinternal.RegisterPersonRegistryServer(s, &server{
		reply: &grpcinternal.RegisterReply{Id: 1},
	})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otel.UnaryClientInterceptor(tracer)),
	)
	if err != nil {
		b.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcinternal.NewPersonRegistryClient(conn)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_, err = client.Register(ctx, &grpcinternal.RegisterRequest{
			Firstname: "Bugs",
			Lastname:  "Bunny",
			Birthdate: &timestamp.Timestamp{Seconds: int64(n)},
			Confirmed: false,
		})

		if err != nil {
			b.Fatalf("call to Register failed: %v", err)
		}
	}
}
