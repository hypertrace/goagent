package bodyattribute // import "github.com/hypertrace/goagent/sdk/instrumentation/bodyattribute"

import (
	"encoding/base64"
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

// SetTruncatedEncodedBodyAttribute is like SetTruncatedBodyAttribute above but also base64 encodes the
// body. This is usually due to non utf8 bytes in the body eg. for multipart/form-data content type.
// The body attribute name has a ".base64" suffix.
func SetTruncatedEncodedBodyAttribute(attrName string, body []byte, bodyMaxSize int, span sdk.Span) {
	bodyLen := len(body)
	if bodyLen == 0 {
		return
	}

	if bodyLen <= bodyMaxSize {
		span.SetAttribute(attrName+".base64", base64.RawStdEncoding.EncodeToString(body))
		return
	}

	truncatedBody := body[:bodyMaxSize]
	span.SetAttribute(fmt.Sprintf("%s.truncated", attrName), true)
	span.SetAttribute(attrName+".base64", base64.RawStdEncoding.EncodeToString(truncatedBody))
}
