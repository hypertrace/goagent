// +build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/tallstoat/pbparser"
)

func getGoPackageName(opts []pbparser.OptionElement) string {
	for _, opt := range opts {
		if opt.Name == "go_package" {
			return opt.Value
		}
	}
	return ""
}

var zeroValues = map[string]string{
	"string": `""`,
	"bool":   `false`,
	"int":    `0`,
}

var fieldNameReplacer = strings.NewReplacer("Http", "HTTP", "Rpc", "RPC")

func toPublicFieldName(name string) string {
	return fieldNameReplacer.Replace(strings.Title(name))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`Usage: generator PROTO_FILE
Parse PROTO_FILE and generate output value objects`)
		return
	}

	file := os.Args[1]
	pf, err := pbparser.ParseFile(file)
	if err != nil {
		fmt.Printf("Unable to parse proto file %q: %v \n", file, err)
		os.Exit(1)
	}

	fqpn := getGoPackageName(pf.Options)
	pn := strings.Split(fqpn, "/")
	c := fmt.Sprintf("package %s\n\n", pn[len(pn)-1])

	for _, m := range pf.Messages {
		// add fields to the structs
		if len(m.Documentation) > 0 {
			c += "// " + m.Documentation + "\n"
		}
		c += fmt.Sprintf("type %s struct {\n", m.Name)
		for _, mf := range m.Fields {
			fieldName := toPublicFieldName(mf.Name)

			if len(mf.Documentation) > 0 {
				c += "    // " + mf.Documentation + "\n"
			}
			if messageType, ok := mf.Type.(pbparser.NamedDataType); ok {
				c += fmt.Sprintf("    %s *%s `json:\"%s,omitempty\"`\n", fieldName, messageType.Name(), mf.Name)
			} else {
				c += fmt.Sprintf("    %s *%s `json:\"%s,omitempty\"`\n", fieldName, mf.Type.Name(), mf.Name)
			}
		}
		c += "}\n\n"

		// adds getters to the structs
		for _, mf := range m.Fields {
			fieldName := toPublicFieldName(mf.Name)

			if messageType, ok := mf.Type.(pbparser.NamedDataType); ok {
				c += fmt.Sprintf("// Get%s returns the %s\n", fieldName, fieldName)
				c += fmt.Sprintf("func (x *%s) Get%s() *%s {\n", m.Name, fieldName, messageType.Name())
				c += fmt.Sprintf("    return x.%s\n", fieldName)
				c += "}\n\n"
			} else {
				c += fmt.Sprintf("// Get%s returns the %s\n", fieldName, mf.Name)
				c += fmt.Sprintf("func (x *%s) Get%s() %s {\n", m.Name, fieldName, mf.Type.Name())
				c += fmt.Sprintf("    if x.%s == nil { return %s }\n", fieldName, zeroValues[mf.Type.Name()])
				c += fmt.Sprintf("    return *x.%s\n", fieldName)
				c += "}\n\n"
			}
		}

		c += fmt.Sprintf("func (x *%s) loadFromEnv(prefix string, defaultValues *%s) {\n", m.Name, m.Name)
		for _, mf := range m.Fields {
			fieldName := toPublicFieldName(mf.Name)
			envPrefix := strings.ToUpper(strcase.ToSnake(mf.Name))
			if namedType, ok := mf.Type.(pbparser.NamedDataType); ok {
				c += fmt.Sprintf("    if x.%s == nil { x.%s = new(%s) }\n", fieldName, fieldName, namedType.Name())
				c += fmt.Sprintf("    x.%s.loadFromEnv(prefix + \"%s_\", defaultValues.%s)\n", fieldName, envPrefix, fieldName)
			} else {
				_type := mf.Type.Name()
				c += fmt.Sprintf("    if x.%s == nil {\n", fieldName)
				c += fmt.Sprintf(
					"        if val, ok := get%sEnv(prefix + \"%s\"); ok {\n",
					strings.Title(_type),
					envPrefix,
				)
				c += fmt.Sprintf("            x.%s = new(%s)\n", fieldName, _type)
				c += fmt.Sprintf("            *x.%s = val\n", fieldName)
				c += fmt.Sprintf("        } else if defaultValues != nil && defaultValues.%s != nil {\n", fieldName)
				c += fmt.Sprintf("            x.%s = new(%s)\n", fieldName, _type)
				c += fmt.Sprintf("            *x.%s = *defaultValues.%s\n", fieldName, fieldName)
				c += fmt.Sprintf("        }\n")
				c += fmt.Sprintf("    }\n\n")
			}
		}
		c += "}\n\n"
	}

	baseFilename := filepath.Base(file)
	outputFile := baseFilename[0 : len(baseFilename)-6] // 6 = len(".proto")

	err = writeToFile(outputFile+".ps.go", []byte(c))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeToFile(filename string, content []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %v", filename, err)
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write into file %q: %v", filename, err)
	}

	return nil
}
