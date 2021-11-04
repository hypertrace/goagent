package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordingDecissionReturnsFalseOnNoContentType(t *testing.T) {
	assert.Equal(t, false, ShouldRecordBodyOfContentType(headerMapAccessor{http.Header{"A": []string{"B"}}}))
}

func TestRecordingDecissionSuccessOnHeaderSet(t *testing.T) {
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
