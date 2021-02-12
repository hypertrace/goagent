package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hypertrace/goagent/sdk"
)

func setAttributesFromHeaders(_type string, headers http.Header, span sdk.Span, filterAttributes map[string](string)) {
	for key, values := range headers {
		if len(values) == 1 {
			setSpanAttribute(
				span,
				filterAttributes,
				fmt.Sprintf("http.%s.header.%s", _type, strings.ToLower(key)),
				values[0],
			)
			continue
		}

		for index, value := range values {
			setSpanAttribute(
				span,
				filterAttributes,
				fmt.Sprintf("http.%s.header.%s[%d]", _type, strings.ToLower(key), index),
				value,
			)
		}
	}
}

func setSpanAttribute(span sdk.Span, filterAttributes map[string](string), key string, value string) {
	span.SetAttribute(key, value)
	if filterAttributes != nil {
		filterAttributes[key] = value
	}
}
