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
	"time"

	"github.com/fatih/structtag"
	"github.com/iancoleman/strcase"
	"graphql-project/core"
	. "graphql-project/generator"
)

type Field struct {
	Name        string
	Flag        string
	Environment string
	Description string
	Type        string
	Default     string
	HasDefault  bool
	Optional    bool
}

// `flag:"port" env:"PORT" desc:"listen port" default:"80" optional:"true"`
type Tags structtag.Tags

func (field *Field) FlagVar() string {
	return core.Uncapitalize(field.Name)
}

func (field *Field) FlagParser() string {
	switch field.Type {
	case "netip.Addr", "[]byte":
		return "StringVar"
	case "time.Duration":
		return "DurationVar"
	case "uint8", "uint16", "uint32":
		return "UintVar"
	case "int8", "int16", "int32":
		return "IntVar"
	case "float32":
		return "Float64Var"
	default:
		return core.Capitalize(field.Type) + "Var"
	}
}

func (field *Field) FlagType() string {
	switch field.Type {
	case "netip.Addr", "[]byte":
		return "string"
	case "uint8", "uint16", "uint32":
		return "uint"
	case "int8", "int16", "int32":
		return "int"
	case "float32":
		return "float64"
	default:
		return field.Type
	}
}

func intDefault(s *string) (def string, hasDef bool, err error) {
	if s == nil {
		def = "0"
		return
	}
	_, err = strconv.ParseInt(*s, 10, 64)
	if err != nil {
		return
	}
	def = *s
	hasDef = true
	return
}

func uintDefault(s *string) (def string, hasDef bool, err error) {
	if s == nil {
		def = "0"
		return
	}
	_, err = strconv.ParseUint(*s, 10, 64)
	if err != nil {
		return
	}
	def = *s
	hasDef = true
	return
}

func floatDefault(s *string) (def string, hasDef bool, err error) {
	if s == nil {
		def = "0"
		return
	}
	_, err = strconv.ParseFloat(*s, 64)
	if err != nil {
		return
	}
	def = *s
	hasDef = true
	return
}

func durationDefault(s *string) (def string, hasDef bool, err error) {
	if s == nil {
		def = "0"
		return
	}
	d, err := time.ParseDuration(*s)
	if err != nil {
		return
	}
	def = strconv.FormatInt(int64(d), 10)
	hasDef = true
	return
}

func boolDefault(s *string) (def string, hasDef bool, err error) {
	if s == nil {
		def = "false"
		return
	}
	_, err = strconv.ParseBool(*s)
	if err != nil {
		return
	}
	def = *s
	hasDef = true
	return
}

func addrDefault(s *string) (def string, hasDef bool, err error) {
	if s == nil {
		def = `""`
		return
	}
	_, err = core.ParseAddress(*s)
	if err != nil {
		return
	}
	def = core.Quote(*s)
	hasDef = true
	return
}

func stringDefault(s *string) (def string, hasDef bool, err error) {
	if s == nil {
		def = `""`
		return
	}
	def = core.Quote(*s)
	hasDef = true
	return
}

func (field *Field) setDefault(s *string) error {
	var err error
	switch field.Type {
	case "int", "int8", "int16", "int32", "int64":
		field.Default, field.HasDefault, err = intDefault(s)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		field.Default, field.HasDefault, err = uintDefault(s)
	case "float32", "float64":
		field.Default, field.HasDefault, err = floatDefault(s)
	case "time.Duration":
		field.Default, field.HasDefault, err = durationDefault(s)
	case "netip.Addr":
		field.Default, field.HasDefault, err = addrDefault(s)
	case "bool":
		field.Default, field.HasDefault, err = boolDefault(s)
	default:
		field.Default, field.HasDefault, err = stringDefault(s)
	}
	if err != nil {
		return fmt.Errorf("%s: default %w", field.Name, err)
	}
	return nil
}

func (field *Field) Accessor() string {
	switch field.Type {
	case "netip.Addr":
		return "Addr"
	case "[]byte":
		return "Bytes"
	case "time.Duration":
		return "Duration"
	case "string":
		return "Str"
	default:
		return core.Capitalize(field.Type)
	}
}

func (field *Field) Exported() bool {
	return core.IsUpperCase(field.Name, 0)
}

func (field *Field) MethodName() string {
	return core.Capitalize(field.Name)
}

func (tags *Tags) Flag() string {
	flag := GetTagName((*structtag.Tags)(tags), "flag")
	if flag == nil {
		return ""
	}
	return *flag
}

func (tags *Tags) Environment() string {
	env := GetTagName((*structtag.Tags)(tags), "env")
	if env == nil {
		return ""
	}
	return *env
}

func (tags *Tags) Description() string {
	desc := GetTagName((*structtag.Tags)(tags), "desc")
	if desc == nil {
		return ""
	}
	return *desc
}

func (tags *Tags) Default() *string {
	return GetTagName((*structtag.Tags)(tags), "default")
}

func (tags *Tags) Optional() bool {
	optional := GetTagName((*structtag.Tags)(tags), "optional")
	if optional == nil {
		return false
	}
	b, err := strconv.ParseBool(*optional)
	if err != nil {
		return false
	}
	return b
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

func getFields(structType *ast.StructType) ([]Field, error) {
	fields := make([]Field, 0, len(structType.Fields.List)+8)
	for _, f := range structType.Fields.List {
		fieldType := getFieldType(f)
		tags := (*Tags)(FieldTags(f.Tag))
		for _, ident := range f.Names {
			field := Field{
				Name:        ident.Name,
				Flag:        tags.Flag(),
				Environment: tags.Environment(),
				Description: tags.Description(),
				Type:        fieldType,
				Optional:    tags.Optional(),
			}
			if field.Flag == "" {
				field.Flag = strcase.ToKebab(field.Name)
			}
			if field.Environment == "" {
				field.Environment = strcase.ToScreamingSnake(field.Name)
			}
			err := field.setDefault(tags.Default())
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		}
	}
	return fields, nil
}

func generate(fileName string, packageName string, types StructTypes) error {

	funcMap := template.FuncMap{
		"toLowerCamel": strcase.ToLowerCamel,
	}
	tmpl, err := template.New("config.tmpl").Funcs(funcMap).ParseFiles("config.tmpl")
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.Grow(64 * 1024)

	params := make(map[string][]Field)
	for typeName, structType := range types {
		if fields, err := getFields(structType); err != nil {
			return err
		} else {
			params[typeName] = fields
		}
	}

	if err := tmpl.Execute(&buffer, params); err != nil {
		return err
	}

	return Write(buffer.Bytes(), fileName[:strings.Index(fileName, ".go")]+"_gen.go")
}
