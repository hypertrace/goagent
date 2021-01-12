package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSourcesPrecedence(t *testing.T) {
	// defines the config file path
	os.Setenv("HT_CONFIG_FILE", "./testdata/config.json")

	// defines the DataCapture.HTTPHeaders.Request = true
	os.Setenv("HT_DATA_CAPTURE_HTTP_HEADERS_REQUEST", "true")

	// defines the DataCapture.HTTPHeaders.Request = false
	os.Setenv("HT_DATA_CAPTURE_HTTP_HEADERS_RESPONSE", "false")

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
	assert.Equal(t, "camel_service", cfg.GetServiceName().GetValue())
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
