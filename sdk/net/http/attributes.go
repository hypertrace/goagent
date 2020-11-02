package http

import (
	"fmt"
	"net/http"

	"github.com/hypertrace/goagent/sdk"
)

func setAttributesFromHeaders(_type string, headers http.Header, span sdk.Span) {
	for key, values := range headers {
		if len(values) == 1 {
			span.SetAttribute(
				fmt.Sprintf("http.%s.header.%s", _type, key),
				values[0],
			)
			continue
		}

		for index, value := range values {
			span.SetAttribute(
				fmt.Sprintf("http.%s.header.%s[%d]", _type, key, index),
				value,
			)
		}
	}
}
