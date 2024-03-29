package grpc // import "github.com/hypertrace/goagent/sdk/instrumentation/google.golang.org/grpc"

import (
	"fmt"

	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/instrumentation/bodyattribute"
)

// setTruncatedBodyAttribute truncates the body and sets the GRPC body as a span attribute.
// When body is being truncated, we also add a second attribute suffixed by `.truncated` to
// make it clear to the user, body has been modified.
func setTruncatedBodyAttribute(_type string, body []byte, bodyMaxSize int, span sdk.Span) {
	bodyattribute.SetTruncatedBodyAttribute(fmt.Sprintf("rpc.%s.body", _type), body, bodyMaxSize, span)
}
