package config

import (
	"github.com/golang/protobuf/jsonpb"
)

// JSONString returns the config as a JSON string
func (c *AgentConfig) JSONString() string {
	m := jsonpb.Marshaler{EmitDefaults: true}
	if content, err := m.MarshalToString(c); err == nil {
		return string(content)
	}

	return ""
}
