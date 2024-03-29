package grpc // import "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// We need a marshaller that does not omit the zero (e.g. 0 or false) to not to lose potentially
// interesting information.
var marshaler = protojson.MarshalOptions{EmitUnpopulated: false}

// MarshalMessageableJSON marshals a value that can be cast as proto.Message into JSON.
func marshalMessageableJSON(messageable interface{}) ([]byte, error) {
	if msg, ok := messageable.(proto.Message); ok {
		return marshaler.Marshal(msg)
	}

	return nil, nil
}
