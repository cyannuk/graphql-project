//go:build generate

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"log"
	"os"
	"strconv"
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
	Type     string
	identity bool
	auto     bool
}

// `db:"id,pk"`
// `gql:"user"`
// `auto:"true"`
type Tags structtag.Tags

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

func (field *Field) DerefType() string {
	return strings.TrimPrefix(field.Type, "*")
}

func (field *Field) NullableType() string {
	switch field.DerefType() {
	case "time.Time":
		return "NullTime"
	case "Date":
		return "NullDate"
	case "int", "uint", "int64", "uint64":
		return "NullBigInt"
	case "byte", "int8", "int16", "int32", "uint8", "uint16", "uint32", "Role":
		return "NullInt"
	case "float32":
		return "NullFloat"
	case "float64":
		return "NullDouble"
	case "bool":
		return "NullBool"
	case "string":
		return "NullString"
	default:
		log.Fatalf("not supported field type %+v", field.Type)
		return ""
	}
}

func (field *Field) IsPtrType() bool {
	return strings.HasPrefix(field.Type, "*")
}

func (field *Field) GoType() string {
	t := field.DerefType()
	switch t {
	case "time.Time":
		return "Time"
	default:
		return t
	}
}

func (field *Field) DbType() string {
	switch field.DerefType() {
	case "int", "uint", "uint64":
		return "int64"
	case "uint32", "Role":
		return "int32"
	case "byte", "int8", "uint8", "uint16":
		return "int16"
	case "time.Time":
		return "Timestamp"
	default:
		return ""
	}
}

// db column name
func (tags *Tags) Column(fieldName string) string {
	dbc := GetTagName((*structtag.Tags)(tags), "db")
	if dbc == nil {
		return strcase.ToLowerCamel(fieldName)
	}
	return *dbc
}

// GQL property name
func (tags *Tags) Property(fieldName string) string {
	gql := GetTagName((*structtag.Tags)(tags), "gql")
	if gql == nil {
		return strcase.ToLowerCamel(fieldName)
	}
	return *gql
}

func (tags *Tags) IsAuto() bool {
	s := GetTagName((*structtag.Tags)(tags), "auto")
	if s == nil {
		return false
	}
	if b, err := strconv.ParseBool(*s); err != nil {
		log.Fatalf("invalid `auto` tag value '%s'", s)
		return false
	} else {
		return b
	}
}

func (tags *Tags) IsPK() bool {
	pk := GetTagOption((*structtag.Tags)(tags), "db", 0)
	return pk != nil && strings.ToLower(*pk) == "pk"
}

func main() {
	if err := Generate(os.Getenv("GOPACKAGE"), ".", os.Getenv("GOFILE"), generate); err != nil {
		log.Fatal(err)
	}
}

func getFields(structType *ast.StructType) []Field {
	fields := make([]Field, 0, len(structType.Fields.List)+8)
	for _, f := range structType.Fields.List {
		tags := (*Tags)(FieldTags(f.Tag))
		for _, ident := range f.Names {
			name := ident.Name
			field := Field{Name: name, Column: tags.Column(name), Property: tags.Property(name), Type: getFieldType(f), auto: tags.IsAuto()}
			if tags.IsPK() || (strings.ToLower(field.Column) == "id") {
				field.identity = true
				field.auto = true
			}
			fields = append(fields, field)
		}
	}
	return fields
}

func generate(fileName string, packageName string, types StructTypes) error {

	funcMap := template.FuncMap{
		"inputs":       inputFields,
		"columns":      columnList,
		"pointers":     fieldValuePtrList,
		"identity":     identity,
		"toLowerCamel": strcase.ToLowerCamel,
		"toSnake":      strcase.ToSnake,
		"plural":       core.Plural,
		"capitalize":   core.Capitalize,
		"join":         core.Join,
	}
	tmpl, err := template.New("entity.tmpl").Funcs(funcMap).ParseFiles("entity.tmpl")
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

func columnList(fields []Field) string {
	b := make([]byte, 0, 256)
	for _, field := range fields {
		b = append(b, '"')
		b = append(b, field.Column...)
		b = append(b, '"')
		b = append(b, ',')
	}
	return gotils.B2S(b)
}

func fieldValuePtrList(objectID string, fields []Field) string {
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

func inputFields(fields []Field) []Field {
	input := make([]Field, 0, len(fields))
	for _, f := range fields {
		if !f.auto {
			input = append(input, f)
		}
	}
	return input
}

func identity(fields []Field) Field {
	for _, field := range fields {
		if field.identity {
			return field
		}
	}
	log.Fatal("identity field not found")
	return Field{}
}
