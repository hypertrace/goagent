package config

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// defaultConfig holds the default config for agent.
var defaultConfig = AgentConfig{
	DataCapture: &DataCapture{
		HttpHeaders: &Message{
			Request:  BoolVal(true),
			Response: BoolVal(true),
		},
		HttpBody: &Message{
			Request:  BoolVal(true),
			Response: BoolVal(true),
		},
		RpcMetadata: &Message{
			Request:  BoolVal(true),
			Response: BoolVal(true),
		},
		RpcBody: &Message{
			Request:  BoolVal(true),
			Response: BoolVal(true),
		},
	},
	Reporting: &Reporting{
		Address: StringVal("localhost"),
		Secure:  BoolVal(false),
	},
}

func BoolVal(val bool) *wrapperspb.BoolValue {
	return &wrapperspb.BoolValue{Value: val}
}

func StringVal(val string) *wrapperspb.StringValue {
	return &wrapperspb.StringValue{Value: val}
}
