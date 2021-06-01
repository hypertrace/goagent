package grpc

import (
	"fmt"

	"github.com/hypertrace/goagent/sdk"
)

// setTruncatedBodyAttribute truncates the body and sets the HTTP body as a span attribute.
// When body is being truncated, we also add a second attribute suffixed by `.truncated` to
// make it clear to the user, body has been modified.
func setTruncatedBodyAttribute(_type string, body []byte, bodyMaxSize int, span sdk.Span) {
	bodyLen := len(body)
	if bodyLen == 0 {
		return
	}

	if bodyLen <= bodyMaxSize {
		span.SetAttribute(fmt.Sprintf("rpc.%s.body", _type), string(body))
		return
	}

	truncatedBody := body[:bodyMaxSize]
	span.SetAttribute(fmt.Sprintf("rpc.%s.body.truncated", _type), true)
	span.SetAttribute(fmt.Sprintf("rpc.%s.body", _type), string(truncatedBody))
}
