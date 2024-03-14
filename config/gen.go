//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/fatih/structtag"
	"github.com/iancoleman/strcase"
	"graphql-project/core"
	. "graphql-project/generator"
)

type Field struct {
	Name        string
	CommandLine string
	Environment string
	Description string
	Default     string
	Type        string
	Bits        int32
	Signed      bool
	Float       bool
}

// `cmd:"port" env:"PORT" desc:"listen port" default:"false"`
type Tags structtag.Tags

func hasType(typeName string, fields []Field) bool {
	for _, field := range fields {
		if field.Type == typeName {
			return true
		}
	}
	return false
}

func (field *Field) VarName() string {
	return core.Uncapitalize(field.Name)
}

func (field *Field) initNumberTraits() {
	switch field.Type {
	case "int8":
		field.Bits = 8
		field.Signed = true
		field.Float = false
	case "int16":
		field.Bits = 16
		field.Signed = true
		field.Float = false
	case "int32":
		field.Bits = 32
		field.Signed = true
		field.Float = false
	case "int64":
		field.Bits = 64
		field.Signed = true
		field.Float = false
	case "uint8":
		field.Bits = 8
		field.Signed = false
		field.Float = false
	case "uint16":
		field.Bits = 16
		field.Signed = false
		field.Float = false
	case "uint32":
		field.Bits = 32
		field.Signed = false
		field.Float = false
	case "uint64":
		field.Bits = 64
		field.Signed = false
		field.Float = false
	case "byte":
		field.Bits = 8
		field.Signed = false
		field.Float = false
	case "int":
		field.Bits = 64
		field.Signed = true
		field.Float = false
	case "uint":
		field.Bits = 64
		field.Signed = false
		field.Float = false
	case "float32":
		field.Bits = 32
		field.Signed = true
		field.Float = true
	case "float64":
		field.Bits = 64
		field.Signed = true
		field.Float = true
	case "time.Duration":
		field.Bits = 64
		field.Signed = true
		field.Float = false
	}
}

func (field *Field) Exported() bool {
	return core.IsUpperCase(field.Name, 0)
}

func (field *Field) MethodName() string {
	return core.Capitalize(field.Name)
}

func (tags *Tags) CommandLine() string {
	return GetTagName((*structtag.Tags)(tags), "cmd", "")
}

func (tags *Tags) Environment() string {
	return GetTagName((*structtag.Tags)(tags), "env", "")
}

func (tags *Tags) Description() string {
	return GetTagName((*structtag.Tags)(tags), "desc", "")
}

func (tags *Tags) Default() string {
	return GetTagName((*structtag.Tags)(tags), "default", "")
}

func main() {
	if err := Generate(os.Getenv("GOPACKAGE"), ".", os.Getenv("GOFILE"), generate); err != nil {
		log.Fatal(err)
	}
}

func getFieldType(field *ast.Field) string {
	if t, ok := field.Type.(*ast.Ident); ok {
		return t.Name
	}
	if t, ok := field.Type.(*ast.SelectorExpr); ok {
		return fmt.Sprintf("%s.%s", t.X, t.Sel)
	}
	if t, ok := field.Type.(*ast.ArrayType); ok {
		return fmt.Sprintf("[]%s", t.Elt)
	}
	if p, ok := field.Type.(*ast.StarExpr); ok {
		if t, ok := p.X.(*ast.Ident); ok {
			return "*" + t.Name
		}
		if t, ok := p.X.(*ast.SelectorExpr); ok {
			return fmt.Sprintf("*%s.%s", t.X, t.Sel)
		}
	}
	log.Fatalf("not supported field type %+v", field.Type)
	return ""
}

func getFields(structType *ast.StructType) []Field {
	fields := make([]Field, 0, len(structType.Fields.List)+8)
	for _, f := range structType.Fields.List {
		fieldType := getFieldType(f)
		tags := (*Tags)(FieldTags(f.Tag))
		for _, ident := range f.Names {
			field := Field{
				Name:        ident.Name,
				CommandLine: tags.CommandLine(),
				Environment: tags.Environment(),
				Description: tags.Description(),
				Default:     tags.Default(),
				Type:        fieldType,
			}
			field.initNumberTraits()
			fields = append(fields, field)
		}
	}
	return fields
}

func generate(fileName string, packageName string, types StructTypes) error {

	funcMap := template.FuncMap{
		"params":       params,
		"toLowerCamel": strcase.ToLowerCamel,
		"hasType":      hasType,
		"join":         core.Join,
	}
	tmpl, err := template.New("config.tmpl").Funcs(funcMap).ParseFiles("config.tmpl")
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.Grow(64 * 10124)

	params := make(map[string][]Field)
	for typeName, structType := range types {
		params[typeName] = getFields(structType)
	}

	if err := tmpl.Execute(&buffer, params); err != nil {
		return err
	}

	return Write(buffer.Bytes(), fileName[:strings.Index(fileName, ".go")]+"_gen.go")
}

type valueHandlerParams struct {
	Field
	ObjectID string
}

func params(field Field, objectID string) valueHandlerParams {
	return valueHandlerParams{field, objectID}
}
