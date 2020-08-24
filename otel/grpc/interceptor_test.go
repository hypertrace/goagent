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
	"github.com/traceableai/goagent/otel/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

var _ PersonRegistryServer = server{}

type server struct {
	*UnimplementedPersonRegistryServer
}

func (server) Register(_ context.Context, _ *RegisterRequest) (*RegisterReply, error) {
	return &RegisterReply{Id: 1}, nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestRegisterPerson(t *testing.T) {
	flusher := internal.InitTracer()

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(NewUnaryServerInterceptor()),
	)
	RegisterPersonRegistryServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := NewPersonRegistryClient(conn)

	_, err = client.Register(ctx, &RegisterRequest{
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

	expectedAssertions := 5
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
