package config

import (
	"os"
	"testing"

	agentconfig "github.com/hypertrace/agent-config/gen/go/v1"
	"github.com/stretchr/testify/assert"
)

// A number of these tests set environment variables.
// When setting an env var, ensure that it is unset
// at the end of the test so as not to impact other tests

func TestSourcesPrecedence(t *testing.T) {
	// defines the config file path
	os.Setenv("HT_CONFIG_FILE", "./testdata/config.json")
	defer os.Unsetenv("HT_CONFIG_FILE")

	// defines the DataCapture.HTTPHeaders.Request = true
	os.Setenv("HT_DATA_CAPTURE_HTTP_HEADERS_REQUEST", "true")
	defer os.Unsetenv("HT_DATA_CAPTURE_HTTP_HEADERS_REQUEST")

	// defines the DataCapture.HTTPHeaders.Request = false
	os.Setenv("HT_DATA_CAPTURE_HTTP_HEADERS_RESPONSE", "false")
	defer os.Unsetenv("HT_DATA_CAPTURE_HTTP_HEADERS_RESPONSE")

	// loads the config
	cfg := Load()
	cfg.DataCapture.RpcMetadata.Response = Bool(false)

	// use defaults
	assert.Equal(t, true, cfg.GetDataCapture().GetHttpBody().GetRequest().GetValue())

	// config file take precedence over defaults
	assert.Equal(t, "http://api.traceable.ai:9411/api/v2/spans", cfg.GetReporting().GetEndpoint().GetValue())

	// env vars take precedence over config file
	assert.Equal(t, false, cfg.GetDataCapture().GetHttpHeaders().GetResponse().GetValue())

	// static value take precedence over config files
	assert.Equal(t, false, cfg.GetDataCapture().GetRpcMetadata().GetResponse().GetValue())
}

func TestCamelYAMLLoadSuccess(t *testing.T) {
	// loads the config
	cfg := LoadFromFile("./testdata/config_camel.yml")

	// config file take precedence over defaults
	assert.Equal(t, "camelService", cfg.GetServiceName().GetValue())
	assert.Equal(t, "http://35.233.143.122:9411/api/v2/spans", cfg.GetReporting().GetEndpoint().GetValue())
	assert.Equal(t, true, cfg.GetDataCapture().GetHttpHeaders().GetRequest().GetValue())
}

func TestSnakeYAMLLoadSuccess(t *testing.T) {
	// loads the config
	cfg := LoadFromFile("./testdata/config_snake.yml")

	// config file take precedence over defaults
	assert.Equal(t, "snake_service", cfg.GetServiceName().GetValue())
	assert.Equal(t, "http://35.233.143.122:9411/api/v2/spans", cfg.GetReporting().GetEndpoint().GetValue())
	assert.Equal(t, true, cfg.GetDataCapture().GetHttpHeaders().GetRequest().GetValue())
}

func TestConfigLoadFromEnvOverridesWithEnv(t *testing.T) {
	cfg := &agentconfig.AgentConfig{
		ServiceName:        agentconfig.String("my_service"),
		PropagationFormats: []agentconfig.PropagationFormat{PropagationFormat_B3, PropagationFormat_TRACECONTEXT},
	}
	assert.Equal(t, "my_service", cfg.GetServiceName().Value)

	os.Setenv("HT_SERVICE_NAME", "my_other_service")
	defer os.Unsetenv("HT_SERVICE_NAME")
	os.Setenv("HT_PROPAGATION_FORMATS", "B3")
	defer os.Unsetenv("HT_PROPAGATION_FORMATS")

	cfg.LoadFromEnv()
	assert.Equal(t, "my_other_service", cfg.GetServiceName().Value)
	assert.Equal(t, 1, len(cfg.GetPropagationFormats()))
	assert.Equal(t, agentconfig.PropagationFormat_B3, cfg.GetPropagationFormats()[0])
}

func TestConfigLoadIsNotOverridenByDefaults(t *testing.T) {
	pf := []agentconfig.PropagationFormat{
		agentconfig.PropagationFormat_B3, agentconfig.PropagationFormat_TRACECONTEXT}

	cfg := &agentconfig.AgentConfig{
		DataCapture: &agentconfig.DataCapture{
			RpcMetadata: &agentconfig.Message{
				Request: agentconfig.Bool(false),
			},
		},
		PropagationFormats: pf,
	}

	assert.Equal(t, false, cfg.DataCapture.RpcMetadata.Request.Value)

	LoadEnv(cfg)
	// we verify here the value isn't overridden by default value (true)
	assert.Equal(t, false, cfg.DataCapture.RpcMetadata.Request.Value)
	// we verify default value is used for undefined value (true)
	assert.Equal(t, true, cfg.DataCapture.RpcMetadata.Response.Value)
	assert.ElementsMatch(t, pf, cfg.GetPropagationFormats())
}
