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
	assert.Equal(t, "api.traceable.ai", cfg.GetReporting().GetAddress().GetValue())

	// env vars take precedence over config file
	assert.Equal(t, false, cfg.GetDataCapture().GetHttpHeaders().GetResponse().GetValue())

	// static value take precedence over config files
	assert.Equal(t, false, cfg.GetDataCapture().GetRpcMetadata().GetResponse().GetValue())

}

func TestYAMLLoadSuccess(t *testing.T) {
	// loads the config
	cfg := LoadFromFile("./testdata/config.yml")

	// config file take precedence over defaults
	assert.Equal(t, "35.233.143.122", cfg.GetReporting().GetAddress().GetValue())
}
