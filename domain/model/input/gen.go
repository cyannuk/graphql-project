//go:build ignore

package main

import (
	"fmt"
	"go/ast"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/fatih/structtag"
	"github.com/iancoleman/strcase"
	gotils "github.com/savsgio/gotils/strconv"

	"graphql-project/core"
	. "graphql-project/generator"
)

type Field struct {
	Name   string
	Column string // db column name
	Type   string
}

type Parameters struct {
	Type   string
	Fields []Field
}

// `dbe:"email"`
type Tags structtag.Tags

func (params *Parameters) HasNullable() bool {
	for _, field := range params.Fields {
		if field.IsNullable() {
			return true
		}
	}
	return false
}

func (field *Field) IsNullable() bool {
	return strings.HasPrefix(field.Type, "model.Null")
}

// db column name
func (tags *Tags) Column(fieldName string) string {
	return GetTagName((*structtag.Tags)(tags), "dbe", strcase.ToLowerCamel(fieldName))
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
	return "string"
}

func getFields(typeName string, structType *ast.StructType) *Parameters {
	params := Parameters{
		Type:   typeName,
		Fields: make([]Field, 0, len(structType.Fields.List)+8),
	}
	for _, f := range structType.Fields.List {
		tags := (*Tags)(FieldTags(f.Tag))
		for _, ident := range f.Names {
			name := ident.Name
			field := Field{Name: name, Column: tags.Column(name), Type: getFieldType(f)}
			params.Fields = append(params.Fields, field)
		}
	}
	return &params
}

func generate(writer io.Writer, packageName string, types StructTypes) error {
	params := struct{ Package string }{packageName}
	if err := headerTemplate.Execute(writer, params); err != nil {
		return err
	}

	for typeName, structType := range types {
		if err := entityImplTemplate.Execute(writer, getFields(typeName, structType)); err != nil {
			return err
		}
	}
	return nil
}

func columnList(fields []Field) string {
	b := make([]byte, 0, 256)
	b = append(b, '`')
	for i, field := range fields {
		if i > 0 {
			b = append(b, ',')
			b = append(b, ' ')
		}
		b = append(b, '"')
		b = append(b, field.Column...)
		b = append(b, '"')
	}
	b = append(b, '`')
	return gotils.B2S(b)
}

func placeholderList(fields []Field) string {
	b := make([]byte, 0, 128)
	b = append(b, '`')
	for i := range fields {
		if i > 0 {
			b = append(b, ',')
			b = append(b, ' ')
		}
		b = append(b, '$')
		b = append(b, core.IntToStr(i+1)...)
	}
	b = append(b, '`')
	return gotils.B2S(b)
}

func valueList(objectID string, fields []Field) string {
	b := make([]byte, 0, 256)
	for i, field := range fields {
		if i > 0 {
			b = append(b, ',')
			b = append(b, ' ')
		}
		b = append(b, objectID...)
		b = append(b, '.')
		b = append(b, field.Name...)
	}
	return gotils.B2S(b)
}

var headerTemplate = template.Must(template.New("header").Parse(`
// Code generated by gen; DO NOT EDIT.
package {{.Package}}

import (
	"graphql-project/domain/model"
)
`))

var entityImplTemplate = template.Must(template.New("entities").
	Funcs(template.FuncMap{
		"columnList":      columnList,
		"placeholderList": placeholderList,
		"valueList":       valueList,
		"toLowerCamel":    strcase.ToLowerCamel,
	}).
	Parse(`
{{$objectID := .Type | toLowerCamel}}
{{- if .HasNullable}}
func ({{$objectID}} *{{.Type}}) Fields() (string, string, []any) {
	f := fields{make([]byte, 0, 128), make([]byte, 0, 64), make([]any, 0, {{len .Fields}})}
{{- range .Fields}}
	{{- if .IsNullable}}
		switch {{$objectID}}.{{.Name}}.State {
		case model.Exists:
			f.addField("{{.Column}}", {{$objectID}}.{{.Name}}.Value)
		case model.Null:
			f.addField("{{.Column}}", nil)
		}
	{{- else}}
		f.addField("{{.Column}}", {{$objectID}}.{{.Name}})
	{{- end}}
{{- end}}
	return f.get()
}
{{- else}}
func ({{$objectID}} *{{.Type}}) Fields() (string, string, []any) {
	return {{columnList .Fields}}, {{placeholderList .Fields}}, []any{ {{valueList $objectID .Fields}} }  
}
{{- end}}
`))
