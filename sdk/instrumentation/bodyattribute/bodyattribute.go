package bodyattribute // import "github.com/hypertrace/goagent/sdk/instrumentation/bodyattribute"

import (
	"encoding/base64"
	"fmt"
	"unicode/utf8"

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
		SetBodyAttribute(attrName, body, false, span)
		return
	}

	truncatedBody := truncateUTF8Bytes(body, bodyMaxSize)

	SetBodyAttribute(attrName, truncatedBody, true, span)
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

	truncatedBody := truncateUTF8Bytes(body, bodyMaxSize)
	SetEncodedBodyAttribute(attrName, truncatedBody, true, span)
}

// SetBodyAttribute sets the body as a span attribute.
// also sets truncated attribute if truncated is true
func SetBodyAttribute(attrName string, body []byte, truncated bool, span sdk.Span) {
	if len(body) == 0 {
		return
	}

	span.SetAttribute(attrName, string(body))
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

	span.SetAttribute(attrName+".base64", base64.RawStdEncoding.EncodeToString(body))
	// if already truncated then set attribute
	if truncated {
		span.SetAttribute(fmt.Sprintf("%s.truncated", attrName), true)
	}
}

// Largely based on:
// https://github.com/jmacd/opentelemetry-go/blob/e8973b75b230246545cdae072a548c83877cba09/sdk/trace/span.go#L358-L375
// Intention here is to ensure that we capture the final parsed rune to prevent splitting multibyte rune in the middle
// If maxBytes - 4 is in middle of rune, that is okay, we still append since max rune size <= 4 so that means there is still
// a partial or full rune between maxBytes - 4 and the end.
// If we encounter a rune that extends beyond the end of our truncation length it will be dropped entirely
func truncateUTF8Bytes(b []byte, maxBytes int) []byte {
	// We subtract 4 as that is the largest possible byte size for single rune
	startIndex := maxBytes - 4
	if startIndex < 0 {
		startIndex = 0
	}

	for idx := startIndex; idx < maxBytes; {
		_, size := utf8.DecodeRune(b[idx:])
		if idx+size > maxBytes {
			// We're past maxBytes with this rune, we will not include this in truncated value
			return b[:idx]
		}
		idx += size
	}

	return b[:maxBytes]
}
