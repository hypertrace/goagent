package grpc

import (
	"testing"

	pb "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc/internal/helloworld"
	"github.com/stretchr/testify/assert"
)

func TestWithProtobufMessageFormatV2(t *testing.T) {
	msg := &pb.HelloReply{Message: "Test Message"}
	data, err := marshalMessageableJSON(msg)
	assert.Nil(t, err)
	assert.Equal(t, "{\"message\":\"Test Message\"}", string(data[:]))
}

func TestWithProtobufMessageFormatV1(t *testing.T) {
	msg := &pb.HelloReply{Message: "Test message"}
	data, err := marshalMessageableJSON(msg)
	assert.Nil(t, err)
	assert.Equal(t, "{\"message\":\"Test message\"}", string(data[:]))
}
