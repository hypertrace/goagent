// +build go1.14

package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyLookupSuccessOnSetForReverseOrder(t *testing.T) {
	tCases := []struct {
		contentTypes []string
		shouldRecord bool
	}{
		{[]string{"charset=utf-8", "application/json"}, true},
	}

	for _, tCase := range tCases {
		h := http.Header{}
		for _, header := range tCase.contentTypes {
			h.Add("Content-Type", header)
		}
		assert.Equal(t, tCase.shouldRecord, shouldRecordBodyOfContentType(h))
	}
}
