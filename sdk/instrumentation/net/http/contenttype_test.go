package http

import (
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordingDecissionReturnsFalseOnNoContentType(t *testing.T) {
	assert.Equal(t, false, ShouldRecordBodyOfContentType(headerMapAccessor{http.Header{"A": []string{"B"}}}))
}

func TestRecordingDecissionSuccessOnHeaderSet(t *testing.T) {
	internalconfig.ResetConfig()
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
}

func TestRecordingDecissionSuccessOnHeaderAdd(t *testing.T) {
	internalconfig.ResetConfig()
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
}

func TestXMLRecordingDecisionSuccessOnHeaderAdd(t *testing.T) {
	cfg := internalconfig.GetConfig()
	cfg.DataCapture.AllowedContentTypes = []*wrapperspb.StringValue{wrapperspb.String("xml")}

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
	internalconfig.ResetConfig()
}
