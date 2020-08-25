package client

import (
	"context"
	"encoding/json"
	"log"
	"net"
	reflect "reflect"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	grpcinternal "github.com/traceableai/goagent/otel/grpc/internal"
	"github.com/traceableai/goagent/otel/internal"
	otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var _ grpcinternal.PersonRegistryServer = server{}

type server struct {
	*grpcinternal.UnimplementedPersonRegistryServer
}

func (server) Register(_ context.Context, _ *grpcinternal.RegisterRequest) (*grpcinternal.RegisterReply, error) {
	return &grpcinternal.RegisterReply{Id: 1}, nil
}

func initListener(s *grpc.Server) func(context.Context, string) (net.Conn, error) {
	const bufSize = 1024 * 1024

	listener := bufconn.Listen(bufSize)
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return bufDialer
}

func TestRegisterPersonSuccess(t *testing.T) {
	_, flusher := internal.InitTracer()

	s := grpc.NewServer()
	defer s.Stop()

	grpcinternal.RegisterPersonRegistryServer(s, &server{})

	dialer := initListener(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(NewUnaryClientInterceptor()),
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

func BenchmarkRequestResponseBodyMarshaling(b *testing.B) {
	internal.InitTracer()

	s := grpc.NewServer()
	defer s.Stop()

	grpcinternal.RegisterPersonRegistryServer(s, &server{})

	dialer := initListener(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(NewUnaryClientInterceptor()),
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

func BenchmarkRequestDefaultInterceptor(b *testing.B) {
	tracer, _ := internal.InitTracer()

	s := grpc.NewServer()
	defer s.Stop()

	grpcinternal.RegisterPersonRegistryServer(s, &server{})

	dialer := initListener(s)

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otel.UnaryClientInterceptor(tracer)),
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
