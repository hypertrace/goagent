package http

import (
	agentconfig "github.com/hypertrace/agent-config/gen/go/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func RestoreDefaults() {
	SetContentTypeAllowList(&agentconfig.AgentConfig{
		DataCapture: &agentconfig.DataCapture{
			AllowedContentTypes: []*wrapperspb.StringValue{wrapperspb.String("json"), wrapperspb.String("x-www-form-urlencoded")},
		},
	})
}

func TestRecordingDecissionReturnsFalseOnNoContentType(t *testing.T) {
	assert.Equal(t, false, ShouldRecordBodyOfContentType(headerMapAccessor{http.Header{"A": []string{"B"}}}))
}

func TestRecordingDecissionSuccessOnHeaderSet(t *testing.T) {
	SetContentTypeAllowList(&agentconfig.AgentConfig{
		DataCapture: &agentconfig.DataCapture{
			AllowedContentTypes: []*wrapperspb.StringValue{wrapperspb.String("json"), wrapperspb.String("x-www-form-urlencoded")},
		},
	})
	tCases := []struct {
		contentType  string
		shouldRecord bool
	}{
		{"text/plain", false},
		{"application/json", true},
		{"Application/JSON", true},
		{"application/json", true},
		{"application/json; charset=utf-8", true},
		{"application/x-www-form-urlencoded", true},
		{"application/vnd.api+json", true},
		{"application/grpc+json", true},
	}

	for _, tCase := range tCases {
		h := http.Header{}
		h.Set("Content-Type", tCase.contentType)
		assert.Equal(t, tCase.shouldRecord, ShouldRecordBodyOfContentType(headerMapAccessor{h}))
	}
	RestoreDefaults()
}

func TestRecordingDecissionSuccessOnHeaderAdd(t *testing.T) {
	SetContentTypeAllowList(&agentconfig.AgentConfig{
		DataCapture: &agentconfig.DataCapture{
			AllowedContentTypes: []*wrapperspb.StringValue{wrapperspb.String("json"), wrapperspb.String("x-www-form-urlencoded")},
		},
	})
	tCases := []struct {
		contentTypes []string
		shouldRecord bool
	}{
		{[]string{"text/plain"}, false},
		{[]string{"application/json"}, true},
		{[]string{"application/json", "charset=utf-8"}, true},
		{[]string{"application/json; charset=utf-8"}, true},
		{[]string{"application/x-www-form-urlencoded"}, true},
		{[]string{"charset=utf-8", "application/json"}, true},
		{[]string{"charset=utf-8", "application/vnd.api+json"}, true},
	}

	for _, tCase := range tCases {
		h := http.Header{}
		for _, header := range tCase.contentTypes {
			h.Add("Content-Type", header)
		}
		assert.Equal(t, tCase.shouldRecord, ShouldRecordBodyOfContentType(headerMapAccessor{h}))
	}
	RestoreDefaults()
}

func TestXMLRecordingDecisionSuccessOnHeaderAdd(t *testing.T) {
	SetContentTypeAllowList(&agentconfig.AgentConfig{
		DataCapture: &agentconfig.DataCapture{
			AllowedContentTypes: []*wrapperspb.StringValue{wrapperspb.String("json"),
				wrapperspb.String("xml")},
		},
	})
	tCases := []struct {
		contentTypes []string
		shouldRecord bool
	}{
		{[]string{"text/xml"}, true},
		{[]string{"application/xml"}, true},
		{[]string{"image/svg+xml"}, true},
		{[]string{"application/xhtml+xml"}, true},
		{[]string{"text/plain"}, false},
	}

	for _, tCase := range tCases {
		h := http.Header{}
		for _, header := range tCase.contentTypes {
			h.Add("Content-Type", header)
		}
		assert.Equal(t, tCase.shouldRecord, ShouldRecordBodyOfContentType(headerMapAccessor{h}))
	}
	RestoreDefaults()
}
