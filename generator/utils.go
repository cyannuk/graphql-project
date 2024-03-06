package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/fatih/structtag"
)

type StructTypes = map[string]*ast.StructType

func findTypes(file *ast.File) StructTypes {
	types := make(StructTypes)
	ast.Inspect(file, func(node ast.Node) bool {
		if typeSpec, ok := node.(*ast.TypeSpec); ok {
			if t, ok := typeSpec.Type.(*ast.StructType); ok {
				types[typeSpec.Name.Name] = t
				return false
			}
		}
		return true
	})
	return types
}

func FieldTags(tag *ast.BasicLit) *structtag.Tags {
	if tag == nil {
		return nil
	}
	if tags, err := structtag.Parse(tag.Value[1 : len(tag.Value)-1]); err != nil {
		return nil
	} else {
		return tags
	}
}

func GetTagName(tags *structtag.Tags, key string, defaultValue string) string {
	if tags != nil {
		if tag, err := tags.Get(key); err != nil {
			return defaultValue
		} else {
			return tag.Name
		}
	}
	return defaultValue
}

func GetTagOption(tags *structtag.Tags, key string, i int, defaultValue string) string {
	if tags != nil {
		if tag, err := tags.Get(key); err != nil {
			return defaultValue
		} else if i < len(tag.Options) {
			return tag.Options[i]
		}
	}
	return defaultValue
}

func Generate(generate func(io.Writer, string, StructTypes) error) error {

	pkgName := os.Getenv("GOPACKAGE")
	srcFile := os.Getenv("GOFILE")

	packages, err := parser.ParseDir(token.NewFileSet(), ".", nil, 0)
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if pkg.Name == pkgName {
			for fileName, pkgFile := range pkg.Files {
				if fileName == srcFile {

					types := findTypes(pkgFile)
					if len(types) == 0 {
						continue
					}

					file, err := os.Create(strings.Replace(fileName, ".go", "_gen.go", 1))
					if err != nil {
						return err
					}

					err = generate(file, pkg.Name, types)
					_ = file.Close()
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
