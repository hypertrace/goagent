package bodyattribute // import "github.com/hypertrace/goagent/sdk/instrumentation/bodyattribute"

import (
	"encoding/base64"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/hypertrace/goagent/sdk"
)

const utf8Replacement = "�"

// SetTruncatedBodyAttribute truncates the body and sets the body as a span attribute.
// When body is being truncated, we also add a second attribute suffixed by `.truncated` to
// make it clear to the user, body has been modified.
func SetTruncatedBodyAttribute(attrName string, body []byte, bodyMaxSize int, span sdk.Span) {
	bodyLen := len(body)
	if bodyLen == 0 {
		return
	}

	if bodyLen <= bodyMaxSize {
		SetBodyAttribute(attrName, body, false, span)
		return
	}

	SetBodyAttribute(attrName, body[:bodyMaxSize], true, span)
}

// SetTruncatedEncodedBodyAttribute is like SetTruncatedBodyAttribute above but also base64 encodes the
// body. This is usually due to non utf8 bytes in the body eg. for multipart/form-data content type.
// The body attribute name has a ".base64" suffix.
func SetTruncatedEncodedBodyAttribute(attrName string, body []byte, bodyMaxSize int, span sdk.Span) {
	bodyLen := len(body)
	if len(body) == 0 {
		return
	}

	if bodyLen <= bodyMaxSize {
		SetEncodedBodyAttribute(attrName, body, false, span)
		return
	}

	SetEncodedBodyAttribute(attrName, body[:bodyMaxSize], true, span)
}

// SetBodyAttribute sets the body as a span attribute.
// also sets truncated attribute if truncated is true
func SetBodyAttribute(attrName string, body []byte, truncated bool, span sdk.Span) {
	if len(body) == 0 {
		return
	}

	bodyStr := string(body)
	if !utf8.ValidString(bodyStr) {
		bodyStr = strings.ToValidUTF8(bodyStr, utf8Replacement)
	}

	span.SetAttribute(attrName, bodyStr)
	// if already truncated then set attribute
	if truncated {
		span.SetAttribute(fmt.Sprintf("%s.truncated", attrName), true)
	}
}

// SetEncodedBodyAttribute is like SetBodyAttribute above but also base64 encodes the
// body. This is usually due to non utf8 bytes in the body eg. for multipart/form-data content type.
// The body attribute name has a ".base64" suffix.
func SetEncodedBodyAttribute(attrName string, body []byte, truncated bool, span sdk.Span) {
	if len(body) == 0 {
		return
	}

	bodyStr := string(body)
	if !utf8.ValidString(bodyStr) {
		bodyStr = strings.ToValidUTF8(bodyStr, utf8Replacement)
	}

	span.SetAttribute(attrName+".base64", base64.RawStdEncoding.EncodeToString([]byte(bodyStr)))
	// if already truncated then set attribute
	if truncated {
		span.SetAttribute(fmt.Sprintf("%s.truncated", attrName), true)
	}
}
