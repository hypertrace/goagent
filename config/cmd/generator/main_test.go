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
	expectedCode, _ := format.Source([]byte(`// PutMyFields sets values in the MyFields map.
	func (x *MyType) PutMyFields(m map[string]string) {
		if len(m) == 0 {
			return
		}
		if x.MyFields == nil {
			x.MyFields = make(map[string]string)
		}
		for k, v := range m {
			if _, ok := x.MyFields[k]; !ok {
				x.MyFields[k] = v
			}
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

	assert.Equal(t, string(expectedCode), string(actualCode))
}

func TestGetSubtypeFromMap(t *testing.T) {
	for _, _type := range []string{"map<string,string>", "map< string, string>"} {
		kType, vType := getSubtypesFromMap(_type)
		assert.Equal(t, "string", kType)
		assert.Equal(t, "string", vType)
	}
}
