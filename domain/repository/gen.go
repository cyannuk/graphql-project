//go:build ignore

package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"graphql-project/core"
	. "graphql-project/generator"
)

func getModelName(model string) string {
	i := strings.IndexByte(model, '.')
	if i < 0 {
		return model
	} else {
		return model[i+1:]
	}
}

func findType(file *ast.File, prefix string, suffix string) string {
	result := ""
	ast.Inspect(file, func(node ast.Node) bool {
		if typeSpec, ok := node.(*ast.TypeSpec); ok {
			typeName := typeSpec.Name.Name
			name := strings.ToLower(typeName)
			if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix) {
				result = typeName
				return false
			}
		}
		return true
	})
	return result
}

func main() {

	model := GetArg(1)
	if model == "" {
		log.Fatal("model not specified")
	}

	srcName := os.Getenv("GOFILE")
	pkgPath := "."

	file, err := parser.ParseFile(token.NewFileSet(), path.Join(pkgPath, srcName), nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	repoTypeName := findType(file, strings.ToLower(getModelName(model)), "repository")
	if repoTypeName == "" {
		log.Fatal("repository type not found")
	}
	err = generate(srcName, model, repoTypeName)
	if err != nil {
		log.Fatal(err)
	}
}

type params struct {
	ModelName string
	ModelType string
	RepoType  string
}

func generate(fileName string, modelTypeName string, repoTypeName string) error {

	funcMap := template.FuncMap{
		"toLowerCamel": strcase.ToLowerCamel,
		"toSnake":      strcase.ToSnake,
		"join":         core.Join,
		"plural":       core.Plural,
	}
	tmpl, err := template.New("repository.tmpl").Funcs(funcMap).ParseFiles("repository.tmpl")
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.Grow(64 * 10124)

	if err := tmpl.Execute(&buffer, &params{getModelName(modelTypeName), modelTypeName, repoTypeName}); err != nil {
		return err
	}

	return Write(buffer.Bytes(), fileName[:strings.Index(fileName, ".go")]+"_gen.go")
}
