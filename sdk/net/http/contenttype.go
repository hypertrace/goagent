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
	var contentTypeValues []string
	for key, values := range h {
		if http.CanonicalHeaderKey(key) == "Content-Type" {
			contentTypeValues = values
		}
	}

	for _, contentTypeValue := range contentTypeValues {
		for _, contentTypeAllowed := range contentTypeAllowList {
			// userland code can set joint headers directly instead of adding
			// them for example
			//
			// ```
			//   header.Set("content-type", "application/json; charset=utf-8")
			// ```
			//
			// instead of
			//
			// ```
			//    header.Add("content-type", "application/json")
			//    header.Add("content-type", "charset=utf-8")
			// ```
			// hence we need to inspect it by using contains.
			if strings.Contains(strings.ToLower(contentTypeValue), contentTypeAllowed) {
				return true
			}
		}
	}
	return false
}
