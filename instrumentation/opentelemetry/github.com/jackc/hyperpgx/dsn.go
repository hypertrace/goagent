package hyperpgx

import (
	"fmt"
	"net/url"
	"strings"
)

// parseDSN parses the connection string provided for postgres driver.
func parseDSN(dsn string) (map[string]string, error) {
	if !strings.HasPrefix(dsn, "postgres://") {
		return nil, fmt.Errorf("invalid postgres connection string as it isn't prefixed with \"postgres://\": %q", dsn)
	}

	// we prepend "http://" to be able to parse the connection string as a URL as it has the same structure.
	parsedURL, err := url.Parse("http://" + dsn[11:])  // 11 = len("postgres://")
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %v", err)
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
