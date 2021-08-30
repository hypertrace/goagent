package hyperpgx // import "github.com/hypertrace/goagent/instrumentation/opentelemetry/github.com/jackc/hyperpgx"

import (
	"fmt"
	"net/url"
)

// parseDSN parses the connection string provided for postgres driver.
func parseDSN(dsn string) (map[string]string, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection %q string: %v", dsn, err)
	}

	if parsedURL.Scheme != "postgres" {
		return nil, fmt.Errorf("invalid postgres connection string %q: it should use \"postgres\" as protocol", dsn)
	}

	connAttrs := map[string]string{}

	if parsedURL.User != nil && parsedURL.User.Username() != "" {
		connAttrs["db.user"] = parsedURL.User.Username()
	} else if queryUser := parsedURL.Query().Get("user"); queryUser != "" {
		connAttrs["db.user"] = queryUser
	}

	if parsedURL.Hostname() != "" {
		connAttrs["net.peer.name"] = parsedURL.Hostname()
	}

	if parsedURL.Port() != "" {
		connAttrs["net.peer.port"] = parsedURL.Port()
	}

	if parsedURL.Path != "/" && parsedURL.Path != "" {
		connAttrs["db.name"] = parsedURL.Path[1:]
	}

	return connAttrs, nil
}
