package grpc

import (
	"context"
	"encoding/json"
	"log"
	"net"
	reflect "reflect"
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
	"google.golang.org/grpc/test/bufconn"
)

var _ grpcinternal.PersonRegistryServer = server{}

type server struct {
	reply        *grpcinternal.RegisterReply
	err          error
	replyHeader  metadata.MD
	replyTrailer metadata.MD
	*grpcinternal.UnimplementedPersonRegistryServer
}

func (s server) Register(ctx context.Context, _ *grpcinternal.RegisterRequest) (*grpcinternal.RegisterReply, error) {
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

func TestServerRegisterPersonSuccess(t *testing.T) {
	_, flusher := internal.InitTracer()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			WrapUnaryServerInterceptor(
				otel.UnaryServerInterceptor(global.TraceProvider().Tracer("ai.traceable")),
			),
		),
	)
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
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcinternal.NewPersonRegistryClient(conn)

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("test_key", "test_value"))
	_, err = client.Register(ctx, &grpcinternal.RegisterRequest{
		Firstname: "Bugs",
		Lastname:  "Bunny",
		Birthdate: &timestamp.Timestamp{Seconds: 1},
		Confirmed: false,
	})
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
	assert.Equal(t, "test_value", attrs.Get("rpc.request.metadata.test_key").AsString())

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

func TestServerRegisterPersonFails(t *testing.T) {
	_, flusher := internal.InitTracer()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			WrapUnaryServerInterceptor(
				otel.UnaryServerInterceptor(global.TraceProvider().Tracer("ai.traceable")),
			),
		),
	)
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

func BenchmarkServerRequestResponseBodyMarshaling(b *testing.B) {
	internal.InitTracer()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			WrapUnaryServerInterceptor(
				otel.UnaryServerInterceptor(global.TraceProvider().Tracer("ai.traceable")),
			),
		),
	)
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

func BenchmarkServerRequestDefaultInterceptor(b *testing.B) {
	tracer, _ := internal.InitTracer()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otel.UnaryServerInterceptor(tracer)),
	)
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
