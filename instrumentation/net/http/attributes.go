package http

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/api/trace"
)

func setAttributesFromHeaders(t string, headers http.Header, span trace.Span) {
	for key, values := range headers {
		if len(values) == 1 {
			span.SetAttribute(
				fmt.Sprintf("http.%s.headers.%s", t, key),
				values[0],
			)
			continue
		}

		for index, value := range values {
			span.SetAttribute(
				fmt.Sprintf("http.%s.headers.%s[%d]", t, key, index),
				value,
			)
		}
	}
}
