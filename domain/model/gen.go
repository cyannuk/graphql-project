//go:build ignore

package main

import (
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
	Name     string
	Column   string // db column name
	Property string // GQL property name
}

type Parameters struct {
	Type     string
	Fields   []Field
	Identity Field
}

// `dbe:"id,pk"`
// `gql:"user"`
type Tags structtag.Tags

func (field *Field) IsID() bool {
	return strings.ToLower(field.Column) == "id"
}

// db column name
func (tags *Tags) Column(fieldName string) string {
	return GetTagName((*structtag.Tags)(tags), "dbe", strcase.ToLowerCamel(fieldName))
}

// GQL property name
func (tags *Tags) Property(fieldName string) string {
	return GetTagName((*structtag.Tags)(tags), "gql", strcase.ToLowerCamel(fieldName))
}

func (tags *Tags) IsPK() bool {
	pk := GetTagOption((*structtag.Tags)(tags), "dbe", 0, "")
	return strings.ToLower(pk) == "pk"
}

func main() {
	if err := Generate(os.Getenv("GOPACKAGE"), ".", os.Getenv("GOFILE"), generate); err != nil {
		log.Fatal(err)
	}
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
			field := Field{Name: name, Column: tags.Column(name), Property: tags.Property(name)}
			params.Fields = append(params.Fields, field)
			if tags.IsPK() || field.IsID() {
				params.Identity = field
			}
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

func fieldList(objectID string, fields []Field) string {
	b := make([]byte, 0, 256)
	for i, field := range fields {
		if i > 0 {
			b = append(b, ',')
			b = append(b, ' ')
		}
		b = append(b, '&')
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
	"graphql-project/interface/model"
)
`))

var entityImplTemplate = template.Must(template.New("entities").
	Funcs(template.FuncMap{
		"columnList":   columnList,
		"fieldList":    fieldList,
		"toLowerCamel": strcase.ToLowerCamel,
		"toSnake":      strcase.ToSnake,
		"plural":       core.Plural,
	}).
	Parse(`
{{$objectID := .Type | toLowerCamel}}
func ({{$objectID}} *{{.Type}}) Table() string {
	return "{{.Type | toSnake | plural}}"
}

func ({{$objectID}} *{{.Type}}) Field(property string) (string, any) {
	switch property {
{{- range .Fields}}
	case "{{.Property}}":
		return "{{.Column}}", &{{$objectID}}.{{.Name}}
{{- end}}
	default:
		return "", nil
	}
}

func ({{$objectID}} *{{.Type}}) Fields() (string, []any) {
	return {{columnList .Fields}}, []any{ {{fieldList $objectID .Fields}} }
}

func ({{$objectID}} *{{.Type}}) Identity() (string, any) {
	return "{{.Identity.Column}}", &{{$objectID}}.{{.Identity.Name}}
}

{{$Types := .Type | plural}}
{{$sliceType := $Types | toLowerCamel}}
{{$slicePtrType := print "p" $sliceType}}

type {{$sliceType}} []{{.Type}}
type {{$slicePtrType}} []*{{.Type}}

{{$slice := $sliceType}}
{{$varName := .Type | toLowerCamel}}

func ({{$slice}} *{{$sliceType}}) New() model.Entity {
	return &{{.Type}}{}
}

func ({{$slice}} *{{$sliceType}}) Add(entity model.Entity) {
	{{$varName}} := entity.(*{{.Type}})
	*{{$slice}} = append(*{{$slice}}, *{{$varName}})
}

func ({{$slice}} *{{$slicePtrType}}) New() model.Entity {
	return &{{.Type}}{}
}

func ({{$slice}} *{{$slicePtrType}}) Add(entity model.Entity) {
	{{$varName}} := *entity.(*{{.Type}})
	*{{$slice}} = append(*{{$slice}}, &{{$varName}})
}

func New{{$Types}}(capacity int) {{$sliceType}} {
	return make([]{{.Type}}, 0, capacity)
}

func NewPtr{{$Types}}(capacity int) {{$slicePtrType}} {
	return make([]*{{.Type}}, 0, capacity)
}
`))
