package config

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// defaultConfig holds the default config values for agent.
var defaultConfig = AgentConfig{
	PropagationFormats: []PropagationFormat{PropagationFormat_TRACECONTEXT},
	DataCapture: &DataCapture{
		HttpHeaders: &Message{
			Request:  Bool(true),
			Response: Bool(true),
		},
		HttpBody: &Message{
			Request:  Bool(true),
			Response: Bool(true),
		},
		RpcMetadata: &Message{
			Request:  Bool(true),
			Response: Bool(true),
		},
		RpcBody: &Message{
			Request:  Bool(true),
			Response: Bool(true),
		},
	},
	Reporting: &Reporting{
		Endpoint: String("localhost:4317"),
		Secure:   Bool(false),
		TraceReporterType: TraceReporterType_OTLP,
	},
}

// Bool wraps the scalar value to be used in the AgentConfig object
func Bool(val bool) *wrapperspb.BoolValue {
	return &wrapperspb.BoolValue{Value: val}
}

// String wraps the scalar value to be used in the AgentConfig object
func String(val string) *wrapperspb.StringValue {
	return &wrapperspb.StringValue{Value: val}
}
