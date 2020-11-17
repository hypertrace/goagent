package hypersql

// Highly inspired in https://github.com/go-sql-driver/mysql/blob/d2e52fca0b/dsn.go

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	errInvalidDSNUnescaped = errors.New("invalid DSN: Unespaced params")
	errInvalidDSNAddr      = errors.New("invalid DSN: network address missing closing bracket")
	errInvalidDSNNoSlash   = errors.New("invalid DSN: missing the slash separating the database name")
)

func parseDSN(dsn string) (map[string]string, error) {
	var (
		user      string
		transport string
		ip        string
		hostport  string
		hostname  string
		port      string
		dbName    string
	)
	// [user[:password]@][transport[(addr)]]/dbname[?param1=value1&paramN=valueN]
	foundSlash := false
	if i := strings.LastIndex(dsn, "/"); i > -1 {
		foundSlash = true
		var j, k int

		// left part is empty if i <= 0
		if i > 0 {
			// [username[:password]@][protocol[(address)]]
			// Find the last '@' in dsn[:i]
			for j = i; j >= 0; j-- {
				if dsn[j] == '@' {
					// username[:password]
					// Find the first ':' in dsn[:j]
					for k = 0; k < j; k++ {
						if dsn[k] == ':' {
							break
						}
					}
					user = dsn[:k]

					break
				}
			}

			// [protocol[(address)]]
			// Find the first '(' in dsn[j+1:i]
			for k = j + 1; k < i; k++ {
				if dsn[k] == '(' {
					// dsn[i-1] must be == ')' if an address is specified
					if dsn[i-1] != ')' {
						if strings.ContainsRune(dsn[k+1:i], ')') {
							return nil, errInvalidDSNUnescaped
						}
						return nil, errInvalidDSNAddr
					}
					hostport = dsn[k+1 : i-1]
					break
				}
			}
			transport = dsn[j+1 : k]
		}

		// dbname[?param1=value1&...&paramN=valueN]
		// Find the first '?' in dsn[i+1:]
		for j = i + 1; j < len(dsn); j++ {
			if dsn[j] == '?' {
				break
			}
		}
		dbName = dsn[i+1 : j]
	}

	if !foundSlash && len(dsn) > 0 {
		return nil, errInvalidDSNNoSlash
	}

	if transport == "" {
		transport = "tcp"
	}

	if hostport == "" {
		switch transport {
		case "tcp":
			ip = "127.0.0.1"
			port = "3306"
		case "unix":
			hostname = "/tmp/mysql.sock"
		default:
			return nil, fmt.Errorf("default addr for network %q unknown", transport)
		}
	} else if transport == "tcp" {
		hostname, ip, port = parseHostport(hostport)
	} else {
		hostname = hostport
	}

	attrs := map[string]string{}
	if user != "" {
		attrs["db.user"] = user
	}

	if dbName != "" {
		attrs["db.name"] = dbName
	}

	switch transport {
	case "tcp":
		attrs["net.transport"] = "IP.TCP"
	case "unix":
		attrs["net.transport"] = "Unix"
	}

	if port != "" {
		attrs["net.peer.port"] = port
	}
	if ip != "" {
		attrs["net.peer.ip"] = ip
	}
	if hostname != "" {
		attrs["net.peer.name"] = hostname
	}

	return attrs, nil
}

func parseHostport(hostport string) (string, string, string) {
	var (
		hostname string
		ip       string
		port     string
	)
	if strings.Count(hostport, ":") > 1 { // presumably ipv6
		if idx := strings.LastIndex(hostport, "]"); idx == -1 { // no brackes, hence no port
			ip = hostport
			port = "3306"
		} else {
			ip = hostport[1:idx]
			port = hostport[idx+2:]
		}
	} else {
		hostportPieces := strings.SplitN(hostport, ":", 2)
		if len(hostportPieces) == 1 {
			port = "3306"
		} else {
			port = hostportPieces[1]
		}

		if parsedIP := net.ParseIP(hostportPieces[0]); parsedIP == nil {
			hostname = hostportPieces[0]
		} else {
			ip = hostportPieces[0]
		}
	}

	return hostname, ip, port
}
