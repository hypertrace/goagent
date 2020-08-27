package internal

import (
	"net/http"
	"strings"
)

var contentTypeAllowList = []string{
	"application/json",
	"application/x-www-form-urlencoded",
}

// ShouldRecordBodyOfContentType checks if the body is meant
// to be recorded based on the content-type and the fact that body is
// not streamed.
func ShouldRecordBodyOfContentType(h http.Header) bool {
	for _, contentType := range contentTypeAllowList {
		// we look for cases like charset=UTF-8; application/json
		for _, value := range h.Values("content-type") {
			// type and subtype are case insensitive
			if strings.ToLower(value) == strings.ToLower(contentType) {
				return true
			}
		}
	}
	return false
}
