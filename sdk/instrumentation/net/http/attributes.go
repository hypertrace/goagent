package http // import "github.com/hypertrace/goagent/sdk/instrumentation/net/http"

import (
	"fmt"
	"log"
	"strings"

	"github.com/hypertrace/goagent/sdk"
)

// SetAttributesFromHeaders set attributes into span from a HeaderAccessor
func SetAttributesFromHeaders(_type string, headers HeaderAccessor, span sdk.Span) {
	err := headers.ForEachHeader(func(key string, values []string) error {
		if len(values) == 1 {
			span.SetAttribute(
				fmt.Sprintf("http.%s.header.%s", _type, strings.ToLower(key)),
				values[0],
			)
			return nil
		}

		for index, value := range values {
			span.SetAttribute(
				fmt.Sprintf("http.%s.header.%s[%d]", _type, strings.ToLower(key), index),
				value,
			)
		}
		return nil
	})

	if err != nil {
		log.Printf("error occurred while setting attributes from headers: %v\n", err)
	}
}
