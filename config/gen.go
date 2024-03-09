//go:build ignore

package main

import (
	"fmt"
	"go/ast"
	"io"
	"log"
	"os"
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
	Type        string
	Bits        int32
	Signed      bool
	Float       bool
}

type Parameters struct {
	Package  string
	Type     string
	Fields   []Field
}

func (params *Parameters) ImportStrConv() bool {
	for _, field := range params.Fields {
		if field.Bits > 0 || field.Type == "bool" {
			return true
		}
	}
	return false
}

func (params *Parameters) ImportBase64() bool {
	for _, field := range params.Fields {
		if field.Type == "[]byte" {
			return true
		}
	}
	return false
}

func (params *Parameters) ImportTime() bool {
	for _, field := range params.Fields {
		if field.Type == "time.Duration" {
			return true
		}
	}
	return false
}

func (params *Parameters) ImportNetIp() bool {
	for _, field := range params.Fields {
		if field.Type == "netip.Addr" {
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

// `cmd:"port" env:"PORT" desc:"listen port"`
type Tags structtag.Tags

func (tags *Tags) CommandLine() string {
	return GetTagName((*structtag.Tags)(tags), "cmd", "")
}

func (tags *Tags) Environment() string {
	return GetTagName((*structtag.Tags)(tags), "env", "")
}

func (tags *Tags) Description() string {
	return GetTagName((*structtag.Tags)(tags), "desc", "")
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

func getFields(packageName string, typeName string, structType *ast.StructType) Parameters {
	params := Parameters{
		Package:  packageName,
		Type:     typeName,
		Fields:   make([]Field, 0, len(structType.Fields.List)+8),
	}
	for _, f := range structType.Fields.List {
		fieldType := getFieldType(f)
		tags := (*Tags)(FieldTags(f.Tag))
		for _, ident := range f.Names {
			ff := Field{
				Name:        ident.Name,
				CommandLine: tags.CommandLine(),
				Environment: tags.Environment(),
				Description: tags.Description(),
				Type:        fieldType,
			}
			ff.initNumberTraits()
			params.Fields = append(params.Fields, ff)
		}
	}
	return params
}

func generate(writer io.Writer, packageName string, types StructTypes) error {

	for typeName, structType := range types {
		parameters := getFields(packageName, typeName, structType)
		if err := configTemplate.Execute(writer, &parameters); err != nil {
			return err
		}
		break
	}
	return nil
}

type valueHandlerParams struct {
	Field
	ObjectID string
}

func params(field Field, objectID string) valueHandlerParams {
	return valueHandlerParams{field, objectID}
}

var configTemplate = template.Must(template.New("config").Funcs(template.FuncMap{"params": params, "toLowerCamel": strcase.ToLowerCamel}).Parse(`// Code generated by gen; DO NOT EDIT.
package {{.Package}}

import (
{{- if .ImportBase64}}
	"encoding/base64"
{{- end}}
{{- if .ImportNetIp}}
	"net/netip"
{{- end}}
{{- if .ImportStrConv}}
	"strconv"
{{- end}}
{{- if .ImportTime}}
	"time"
{{- end}}
	"bufio"
	"errors"
	"flag"
	"os"
	"strings"
	"graphql-project/core"
)

var (
{{- range .Fields}}
	{{.VarName}} string
{{- end}}
)

func init() {
{{- range .Fields}}
	flag.StringVar(&{{.VarName}}, "{{.CommandLine}}", "", "{{.Description}}")
{{- end}}
	flag.Parse()
}

{{$objectID := .Type | toLowerCamel}}

{{- range .Fields}}
	{{if not .Exported}}
		func ({{$objectID}} *{{$.Type}}) {{.MethodName}}() {{.Type}} {
			return {{$objectID}}.{{.Name}}
		}
	{{end}}
{{- end}}

func loadDotEnvFile(fileName string) (map[string]string, error) {
	if file, err := os.OpenFile(fileName, os.O_RDONLY, 0); err != nil {
		return nil, err
	} else {
		defer file.Close()

		stringBuffer := make([]string, 0, 128)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			stringBuffer = append(stringBuffer, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		values := make(map[string]string)
		for _, str := range stringBuffer {
			if core.StartWith(str, '#') {
				continue
			}
			if i := strings.IndexByte(str, '='); i > 0 {
				param := strings.TrimSpace(str[:i])
				if param != "" {
					value := str[i+1:]
					if i := strings.IndexByte(value, '#'); i >= 0 {
						value = value[:i]
					}
					value = core.TrimQuotes(strings.TrimSpace(value))
					if value != "" {
						values[param] = value
					}
				}
			}
		}

		return values, nil
	}
}

{{- define "valueHandler"}}
	{{- if gt .Bits 0}}
		{{- if .Float}}
			if v, err := strconv.ParseFloat(s, {{.Bits}}); err != nil {
		{{- else if .Signed}}
			if v, err := strconv.ParseInt(s, 10, {{.Bits}}); err != nil {
		{{- else}}
			if v, err := strconv.ParseUint(s, 10, {{.Bits}}); err != nil {
		{{- end}}
				return err
			} else {
				{{.ObjectID}}.{{.Name}} = {{.Type}}(v)
			}
	{{- else if eq .Type "string"}}
		{{.ObjectID}}.{{.Name}} = s
	{{- else if eq .Type "[]byte"}}
		if v, err := base64.StdEncoding.DecodeString(s); err != nil {
			return err
		} else {
			{{.ObjectID}}.{{.Name}} = v
		}
	{{- else if eq .Type "bool"}}
		s = strings.ToLower(s)
		var v bool
		if s == "0" || s == "n" {
			v = false
		} else if s == "1" || s == "y" {
			v = true
		} else if b, err := strconv.ParseBool(s); err != nil {
			return err
		} else {
			v = b
		}
		{{.ObjectID}}.{{.Name}} = v
	{{- else if eq .Type "netip.Addr"}}
		if v, err := netip.ParseAddr(s); err != nil {
			return err
		} else {
			{{.ObjectID}}.{{.Name}} = v
		}
	{{- end}}
	exists["{{.Name}}"] = true
{{- end}}

func ({{$objectID}} *{{.Type}}) loadEnv(exists map[string]bool) error {
{{- range .Fields}}
	if s, ok := os.LookupEnv("{{.Environment}}"); !ok || s == "" {
		return nil
	} else {
		{{- template "valueHandler" params . $objectID}}
	}
{{- end}}
	return nil
}

func ({{$objectID}} *{{.Type}}) loadDotEnv(exists map[string]bool) error {
	if values, err := loadDotEnvFile(".env"); err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	} else {
{{- range .Fields}}
		if s, ok := values["{{.Environment}}"]; ok {
			if s == "" {
				return errors.New("empty configuration parameter: {{.Name}}")
			} else {
				{{- template "valueHandler" params . $objectID}}
			}
		}
{{- end}}
	}
	return nil
}

func ({{$objectID}} *{{.Type}}) loadFlags(exists map[string]bool) error {
	var s string
{{- range .Fields}}
	s = {{.VarName}}
	if s != "" {
		{{- template "valueHandler" params . $objectID}}
	}
{{- end}}
	return nil
}

func ({{$objectID}} *{{.Type}}) Load() error {
	exists := make(map[string]bool)
	if err := {{$objectID}}.loadEnv(exists); err != nil {
		return err
	}
	if err := config.loadDotEnv(exists); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if err := {{$objectID}}.loadFlags(exists); err != nil {
		return err
	}
{{- range .Fields}}
	if v, ok := exists["{{.Name}}"]; !ok || !v {
		return errors.New("no configuration parameter: {{.Name}}")
	}
{{- end}}
	return nil
}
`))
