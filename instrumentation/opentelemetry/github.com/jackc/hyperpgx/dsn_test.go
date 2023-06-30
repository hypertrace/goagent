package hyperpgx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDSN(t *testing.T) {
	type testCase struct {
		inputDSN      string
		expectedAttrs map[string]string
	}

	// Test cases have been taken from https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	tCases := []testCase{
		{inputDSN: "postgres://", expectedAttrs: map[string]string{}},
		{
			inputDSN: "postgres://localhost",
			expectedAttrs: map[string]string{
				"net.peer.name": "localhost",
			},
		},
		{
			inputDSN: "postgres://localhost:5432",
			expectedAttrs: map[string]string{
				"net.peer.name": "localhost",
				"net.peer.port": "5432",
			},
		},
		{
			inputDSN: "postgres://localhost/mydb",
			expectedAttrs: map[string]string{
				"net.peer.name": "localhost",
				"db.name":       "mydb",
			},
		},
		{
			inputDSN: "postgres://user@localhost",
			expectedAttrs: map[string]string{
				"net.peer.name": "localhost",
				"db.user":       "user",
			},
		},
		{
			inputDSN: "postgres://user:secret@localhost",
			expectedAttrs: map[string]string{
				"net.peer.name": "localhost",
				"db.user":       "user",
			},
		},
		{
			inputDSN: "postgres://other@localhost/otherdb?connect_timeout=10&application_name=myapp",
			expectedAttrs: map[string]string{
				"net.peer.name": "localhost",
				"db.user":       "other",
				"db.name":       "otherdb",
			},
		},
		{
			inputDSN: "postgres://localhost/mydb?user=other&password=secret",
			expectedAttrs: map[string]string{
				"net.peer.name": "localhost",
				"db.user":       "other",
				"db.name":       "mydb",
			},
		},
	}

	for _, tCase := range tCases {
		t.Run(tCase.inputDSN[11:], func(t *testing.T) { // 11 = len("postgres://")
			connAttrs, err := parseDSN(tCase.inputDSN)
			require.NoError(t, err)
			assert.EqualValues(t, tCase.expectedAttrs, connAttrs)
		})
	}
}
