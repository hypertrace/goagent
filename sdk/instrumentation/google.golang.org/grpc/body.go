package grpc // import "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"

import (
	"fmt"

	"github.com/hypertrace/goagent/sdk"
)

// setTruncatedBodyAttribute truncates the body and sets the GRPC body as a span attribute.
// If the GRPC is larger than this, ignore it entirely. If the body fits into `max_processing_size`,
// decode the body, and then pass handle the resulting body.
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

	span.SetAttribute(fmt.Sprintf("rpc.%s.body.truncated", _type), true)
}
