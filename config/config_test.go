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

	// defines the DataCapture.HTTPHeaders.Request = true
	os.Setenv("HT_DATA_CAPTURE_HTTP_HEADERS_RESPONSE", "true")

	// loads the config
	cfg := Load()

	// use defaults
	assert.Equal(t, false, cfg.GetDataCapture().GetHTTPBody().GetResponse())

	// config file take precedence over defaults
	assert.Equal(t, "api.traceable.ai", cfg.GetReporting().GetAddress())

	// env vars take precedence over config file
	assert.Equal(t, true, cfg.GetDataCapture().GetHTTPHeaders().GetRequest())

	// static value take precedence over config files
	cfg.DataCapture.HTTPHeaders.Response = BoolVal(false)
}

func TestYAMLLoadSuccess(t *testing.T) {
	// loads the config
	cfg := LoadFromFile("./testdata/config.yml")

	// config file take precedence over defaults
	assert.Equal(t, "35.233.143.122", cfg.GetReporting().GetAddress())
}
