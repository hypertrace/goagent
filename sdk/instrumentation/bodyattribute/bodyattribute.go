package bodyattribute // import "github.com/hypertrace/goagent/sdk/instrumentation/bodyattribute"

import (
	"fmt"

	"github.com/hypertrace/goagent/sdk"
)

// SetTruncatedBodyAttribute truncates the body and sets the body as a span attribute.
// When body is being truncated, we also add a second attribute suffixed by `.truncated` to
// make it clear to the user, body has been modified.
func SetTruncatedBodyAttribute(attrName string, body []byte, bodyMaxSize int, span sdk.Span) {
	bodyLen := len(body)
	if bodyLen == 0 {
		return
	}

	if bodyLen <= bodyMaxSize {
		span.SetAttribute(attrName, string(body))
		return
	}

	truncatedBody := body[:bodyMaxSize]
	span.SetAttribute(fmt.Sprintf("%s.truncated", attrName), true)
	span.SetAttribute(attrName, string(truncatedBody))
}
