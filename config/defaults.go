package config // import "github.com/hypertrace/goagent/config"

import (
	agentconfig "github.com/hypertrace/agent-config/gen/go/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// defaultConfig holds the default config values for agent.
var defaultConfig = agentconfig.AgentConfig{
	Enabled:            agentconfig.Bool(true),
	PropagationFormats: []agentconfig.PropagationFormat{agentconfig.PropagationFormat_TRACECONTEXT},
	DataCapture: &agentconfig.DataCapture{
		HttpHeaders: &agentconfig.Message{
			Request:  agentconfig.Bool(true),
			Response: agentconfig.Bool(true),
		},
		HttpBody: &agentconfig.Message{
			Request:  agentconfig.Bool(true),
			Response: agentconfig.Bool(true),
		},
		RpcMetadata: &agentconfig.Message{
			Request:  agentconfig.Bool(true),
			Response: agentconfig.Bool(true),
		},
		RpcBody: &agentconfig.Message{
			Request:  agentconfig.Bool(true),
			Response: agentconfig.Bool(true),
		},
		BodyMaxSizeBytes:           agentconfig.Int32(131072),
		BodyMaxProcessingSizeBytes: agentconfig.Int32(1048576),
		AllowedContentTypes: []*wrapperspb.StringValue{wrapperspb.String("json"),
			wrapperspb.String("x-www-form-urlencoded")},
	},
	Reporting: &agentconfig.Reporting{
		Endpoint:                agentconfig.String("http://localhost:9411/api/v2/spans"),
		Secure:                  agentconfig.Bool(false),
		TraceReporterType:       agentconfig.TraceReporterType_ZIPKIN,
		CertFile:                agentconfig.String(""),
		EnableGrpcLoadbalancing: agentconfig.Bool(true),
	},
}
