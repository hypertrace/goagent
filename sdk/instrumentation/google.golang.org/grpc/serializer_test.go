package grpc

import (
	"testing"

	pb "github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc/examples/helloworld"
	"github.com/stretchr/testify/assert"
)

func TestWithProtobufMessageFormatV2(t *testing.T){
	msg := &pb.HelloReply{Message: "Test Message"}
	data, err := marshalMessageableJSON(msg)
	assert.Nil(t, err)
	assert.Equal(t, "{\"message\":\"Test Message\"}", string(data[:]))
}

func TestWithProtobufMessageFormatV1(t *testing.T){
	msg := &pb.HelloResponse{Name: "Test Name"}
	data, err := marshalMessageableJSON(msg)
	assert.Nil(t, err)
	assert.Equal(t, "{\"name\":\"Test Name\"}", string(data[:]))
}