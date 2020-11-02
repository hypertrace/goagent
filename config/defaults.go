package config

// defaultConfig holds the default config for agent.
var defaultConfig = AgentConfig{
	DataCapture: &DataCapture{
		HTTPHeaders: &Message{
			Request:  BoolVal(false),
			Response: BoolVal(false),
		},
		HTTPBody: &Message{
			Request:  BoolVal(false),
			Response: BoolVal(false),
		},
		RPCMetadata: &Message{
			Request:  BoolVal(false),
			Response: BoolVal(false),
		},
		RPCBody: &Message{
			Request:  BoolVal(false),
			Response: BoolVal(false),
		},
	},
	Reporting: &Reporting{
		Address: StringVal("localhost"),
		Secure:  BoolVal(false),
	},
}

// StringVal returns the pointer value from a string
func StringVal(s string) *string {
	return &s
}

// BoolVal returns the pointer value from a string
func BoolVal(b bool) *bool {
	return &b
}
