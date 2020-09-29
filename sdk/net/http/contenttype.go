package http

import (
	"net/http"
	"strings"
)

// contentTypeAllowList is the list of allowed content types in lowercase
var contentTypeAllowList = []string{
	"application/json",
	"application/x-www-form-urlencoded",
}

// shouldRecordBodyOfContentType checks if the body is meant
// to be recorded based on the content-type and the fact that body is
// not streamed.
func shouldRecordBodyOfContentType(h http.Header) bool {
	for _, contentType := range contentTypeAllowList {
		if strings.ToLower(h.Get("content-type")) == contentType {
			return true
		}
	}
	return false
}
