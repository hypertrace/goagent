package http // import "github.com/hypertrace/goagent/sdk/instrumentation/net/http"

import (
	agentconfig "github.com/hypertrace/agent-config/gen/go/v1"
	"strings"
)

// contentTypeAllowList is the list of allowed "stems" of content types in lowercase.
// Example headers that will match this include:
// - application/json
// - application/vnd.api+json
// - application/x-www-form-urlencoded'
var contentTypeAllowListLowerCase = []string{
	"json", "x-www-form-urlencoded",
}

func SetContentTypeAllowList(config *agentconfig.AgentConfig) {
	contentTypeAllowListLowerCase = []string{}
	for _, contentType := range config.DataCapture.GetAllowedContentTypes() {
		contentTypeAllowListLowerCase = append(contentTypeAllowListLowerCase, strings.ToLower(contentType.GetValue()))
	}
}

// ShouldRecordBodyOfContentType checks if the body is meant
// to be recorded based on the content-type and the fact that body is
// not streamed.
func ShouldRecordBodyOfContentType(h HeaderAccessor) bool {
	var contentTypeValues = h.Lookup("Content-Type") // "Content-Type" is the canonical key

	// we iterate all the values as userland code add the headers in the inverse order,
	// e.g.
	// ```
	//    header.Add("content-type", "charset=utf-8")
	//    header.Add("content-type", "application/json")
	// ```
	for _, contentTypeValue := range contentTypeValues {
		for _, contentTypeAllowedLowerCase := range contentTypeAllowListLowerCase {
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
			if strings.Contains(strings.ToLower(contentTypeValue), contentTypeAllowedLowerCase) {
				return true
			}
		}
	}
	return false
}
