package internal

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// We need a marshaller that does not omit the zero (e.g. 0 or false) to not to lose pontetially
// intesting information.
var marshaler = protojson.MarshalOptions{EmitUnpopulated: true}

// MarshalMessageableJSON marshals a value that an be cast as proto.Message into JSON.
func MarshalMessageableJSON(messageable interface{}) ([]byte, error) {
	if m, ok := messageable.(proto.Message); ok {
		return marshaler.Marshal(m)
	}

	return nil, nil
}
