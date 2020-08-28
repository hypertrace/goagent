package grpc

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	grpcinternal "github.com/traceableai/goagent/otel/grpc/internal"
	"github.com/traceableai/goagent/otel/internal"
	otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"google.golang.org/grpc"
)

func TestServerRegisterPersonSuccess(t *testing.T) {
	_, flusher := internal.InitTracer()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(NewUnaryServerInterceptor()),
	)
	defer s.Stop()

	grpcinternal.RegisterPersonRegistryServer(s, &server{})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcinternal.NewPersonRegistryClient(conn)

	_, err = client.Register(ctx, &grpcinternal.RegisterRequest{
		Firstname: "Bugs",
		Lastname:  "Bunny",
		Birthdate: &timestamp.Timestamp{Seconds: 1},
		Confirmed: false,
	})
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	spans := flusher()
	assert.Equal(t, 1, len(spans))

	span := spans[0]
	assert.Equal(t, "helloworld.PersonRegistry/Register", span.Name)

	expectedAssertions := 5 // one per each tag
	for _, kv := range span.Attributes {
		switch kv.Key {
		case "rpc.system":
			assert.Equal(t, "grpc", kv.Value.AsString())
		case "rpc.service":
			assert.Equal(t, "helloworld.PersonRegistry", kv.Value.AsString())
		case "rpc.method":
			assert.Equal(t, "Register", kv.Value.AsString())
		case "grpc.request.body":
			expectedBody := "{\"firstname\":\"Bugs\",\"lastname\":\"Bunny\",\"birthdate\":\"1970-01-01T00:00:01Z\",\"confirmed\":false}"
			if ok, err := jsonEqual(expectedBody, kv.Value.AsString()); err == nil {
				assert.True(t, ok)
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		case "grpc.response.body":
			expectedBody := "{\"id\":\"1\"}"
			if ok, err := jsonEqual(expectedBody, kv.Value.AsString()); err == nil {
				assert.True(t, ok)
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		default:
			t.Errorf("unexpected attribute %s", kv.Key)
		}
		expectedAssertions = expectedAssertions - 1
	}

	if expectedAssertions > 0 {
		t.Errorf("unexpected number of assertions, missing %d", expectedAssertions)
	}
}

func BenchmarkServerRequestResponseBodyMarshaling(b *testing.B) {
	internal.InitTracer()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(NewUnaryServerInterceptor()),
	)
	defer s.Stop()

	grpcinternal.RegisterPersonRegistryServer(s, &server{})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		b.Fatalf("Failed to dial bufnet: %v", err)
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
			b.Fatalf("Registration failed: %v", err)
		}
	}
}

func BenchmarkServerRequestDefaultInterceptor(b *testing.B) {
	tracer, _ := internal.InitTracer()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otel.UnaryServerInterceptor(tracer)),
	)
	defer s.Stop()

	grpcinternal.RegisterPersonRegistryServer(s, &server{})

	dialer := createDialer(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		b.Fatalf("Failed to dial bufnet: %v", err)
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
			b.Fatalf("Registration failed: %v", err)
		}
	}
}
