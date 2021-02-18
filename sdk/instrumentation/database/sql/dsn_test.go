package sql

import (
	"reflect"
	"testing"
)

func TestParseValidDSN(t *testing.T) {
	tCases := []struct {
		dsn           string
		expectedAttrs map[string]string
	}{{
		"username:password@protocol(address)/dbname?param=value",
		map[string]string{"db.user": "username", "net.peer.name": "address", "db.name": "dbname"},
	}, {
		"user@unix(/path/to/socket)/dbname?charset=utf8",
		map[string]string{"db.user": "user", "net.peer.name": "/path/to/socket", "net.transport": "Unix", "db.name": "dbname"},
	}, {
		"user:password@tcp(localhost:5555)/dbname?charset=utf8",
		map[string]string{"db.user": "user", "net.transport": "IP.TCP", "net.peer.name": "localhost", "net.peer.port": "5555", "db.name": "dbname"},
	}, {
		"user:password@/dbname?allowNativePasswords=false",
		map[string]string{"db.user": "user", "net.peer.ip": "127.0.0.1", "net.peer.port": "3306", "db.name": "dbname", "net.transport": "IP.TCP"},
	}, {
		"user:p@ss(word)@tcp([de:ad:be:ef::ca:fe]:80)/dbname?loc=Local",
		map[string]string{"db.user": "user", "net.transport": "IP.TCP", "net.peer.ip": "de:ad:be:ef::ca:fe", "net.peer.port": "80", "db.name": "dbname"},
	}, {
		"/dbname",
		map[string]string{"net.transport": "IP.TCP", "net.peer.ip": "127.0.0.1", "net.peer.port": "3306", "db.name": "dbname"},
	}, {
		"@/",
		map[string]string{"net.transport": "IP.TCP", "net.peer.ip": "127.0.0.1", "net.peer.port": "3306"},
	}, {
		"/",
		map[string]string{"net.transport": "IP.TCP", "net.peer.ip": "127.0.0.1", "net.peer.port": "3306"},
	}, {
		"",
		map[string]string{"net.transport": "IP.TCP", "net.peer.ip": "127.0.0.1", "net.peer.port": "3306"},
	}, {
		"user:p@/ssword@/",
		map[string]string{"db.user": "user", "net.transport": "IP.TCP", "net.peer.ip": "127.0.0.1", "net.peer.port": "3306"},
	}, {
		"unix/?arg=%2Fsome%2Fpath.ext",
		map[string]string{"net.peer.name": "/tmp/mysql.sock", "net.transport": "Unix"},
	}, {
		"tcp(127.0.0.1)/dbname",
		map[string]string{"net.transport": "IP.TCP", "net.peer.ip": "127.0.0.1", "net.peer.port": "3306", "db.name": "dbname"},
	}, {
		"tcp(de:ad:be:ef::ca:fe)/dbname",
		map[string]string{"net.transport": "IP.TCP", "net.peer.ip": "de:ad:be:ef::ca:fe", "net.peer.port": "3306", "db.name": "dbname"},
	},
	}

	for _, tCase := range tCases {
		attrs, err := parseDSN(tCase.dsn)
		if err != nil {
			t.Error("unexpected parsed attributes")
		}

		if !reflect.DeepEqual(attrs, tCase.expectedAttrs) {
			t.Errorf("\nDSN: %q mismatch:\ngot  %+v\nwant %+v", tCase.dsn, attrs, tCase.expectedAttrs)
		}
	}
}

func TestParseInvalidDSN(t *testing.T) {
	var invalidDSNs = []string{
		"@net(addr/",                  // no closing brace
		"@tcp(/",                      // no closing brace
		"tcp(/",                       // no closing brace
		"(/",                          // no closing brace
		"net(addr)//",                 // unescaped
		"User:pass@tcp(1.2.3.4:3306)", // no trailing slash
		"net()/",                      // unknown default addr
	}

	for _, invalidDSN := range invalidDSNs {
		if _, err := parseDSN(invalidDSN); err == nil {
			t.Errorf("invalid DSN: expected error %q", invalidDSN)
		}
	}
}
