package http // import "github.com/hypertrace/goagent/sdk/instrumentation/net/http"

import (
	"fmt"

	"github.com/hypertrace/goagent/sdk"
	"github.com/hypertrace/goagent/sdk/instrumentation/bodyattribute"
)

// setTruncatedBodyAttribute truncates the body and sets the HTTP body as a span attribute.
// When body is being truncated, we also add a second attribute suffixed by `.truncated` to
// make it clear to the user, body has been modified. Also if base64Encode == true, we base64
// encode the body and append the suffix ".base64" to the attribute name. We base64 encode in
// in case there are non utf8 bytes in the body eg. for binary files in multipart/form-data
// content-type.
func setTruncatedBodyAttribute(_type string, body []byte, bodyMaxSize int, span sdk.Span, base64Encode bool) {
	if base64Encode {
		bodyattribute.SetTruncatedEncodedBodyAttribute(fmt.Sprintf("http.%s.body", _type), body, bodyMaxSize, span)
	} else {
		bodyattribute.SetTruncatedBodyAttribute(fmt.Sprintf("http.%s.body", _type), body, bodyMaxSize, span)
	}
}
