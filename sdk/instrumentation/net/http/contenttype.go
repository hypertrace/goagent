package http // import "github.com/hypertrace/goagent/sdk/instrumentation/net/http"

import (
	internalconfig "github.com/hypertrace/goagent/sdk/internal/config"
	"strings"
)

const contentTypeHeaderKey string = "Content-Type"

// ShouldRecordBodyOfContentType checks if the body is meant
// to be recorded based on the content-type and the fact that body is
// not streamed.
func ShouldRecordBodyOfContentType(h HeaderAccessor) bool {
	var contentTypeValues = h.Lookup(contentTypeHeaderKey) // "Content-Type" is the canonical key
	cfg := internalconfig.GetConfig().GetDataCapture()
	if cfg == nil || cfg.GetAllowedContentTypes() == nil {
		return false
	}
	// we iterate all the values as userland code add the headers in the inverse order,
	// e.g.
	// ```
	//    header.Add("content-type", "charset=utf-8")
	//    header.Add("content-type", "application/json")
	// ```
	for _, contentTypeValue := range contentTypeValues {
		for _, contentTypeAllowed := range cfg.GetAllowedContentTypes() {
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
			if strings.Contains(strings.ToLower(contentTypeValue), strings.ToLower(contentTypeAllowed.Value)) {
				return true
			}
		}
	}
	return false
}

// HasMultiPartFormDataContentTypeHeader returns true if the Content-Type header is
// multipart/form-data. false otherwise.
func HasMultiPartFormDataContentTypeHeader(h HeaderAccessor) bool {
	var contentTypeValues = h.Lookup(contentTypeHeaderKey)
	for _, contentTypeValue := range contentTypeValues {
		if strings.Contains(strings.ToLower(contentTypeValue), "multipart/form-data") {
			return true
		}
	}
	return false
}
