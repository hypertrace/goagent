package internal

import "net/http"

var contentTypeAllowList = []string{
	"application/json",
	"application/x-www-form-urlencoded",
}

// IsContentTypeInAllowList checks if the body is meant
// to be recorded based on the content-type and the fact that body is
// not streamed.
func IsContentTypeInAllowList(h http.Header) bool {
	for _, contentType := range contentTypeAllowList {
		// we look for cases like charset=UTF-8; application/json
		for _, value := range h.Values("content-type") {
			if value == contentType {
				return true
			}
		}
	}
	return false
}
