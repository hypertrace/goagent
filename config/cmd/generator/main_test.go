package main

import (
	"go/format"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tallstoat/pbparser"
)

type MapDataType struct {
}

func (dt MapDataType) Name() string {
	return "map<string,string>"
}

func (dt MapDataType) Category() pbparser.DataTypeCategory {
	return 99
}

func TestSetMapFieldSetter(t *testing.T) {
	expectedCode, _ := format.Source([]byte(`// SetMyFields sets values in the MyFields map.
	func (x *MyType) SetMyFields(m map[string]string) {
		for k, v := range m {
			x.MyFields[k] = v
		}
	}
`))

	c := ""
	m := pbparser.MessageElement{
		Name: "MyType",
	}

	mf := pbparser.FieldElement{
		Name: "my_fields",
		Type: MapDataType{},
	}

	addMapFieldSetter(&c, m, mf)

	actualCode, err := format.Source([]byte(c))
	require.NoError(t, err)

	assert.Equal(t, expectedCode, actualCode)
}
