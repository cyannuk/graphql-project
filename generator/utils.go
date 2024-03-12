package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/fatih/structtag"

	"graphql-project/core"
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
	if tags, err := structtag.Parse(core.TrimQuotes(tag.Value)); err != nil {
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

func Generate(pkgName string, pkgPath string, srcName string, generate func(io.Writer, string, StructTypes) error) error {

	packages, err := parser.ParseDir(token.NewFileSet(), pkgPath, nil, 0)
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if pkg.Name == pkgName {
			for fileName, pkgFile := range pkg.Files {
				if fileName == srcName {

					types := findTypes(pkgFile)
					if len(types) == 0 {
						continue
					}

					var buffer bytes.Buffer
					buffer.Grow(64 * 10124)
					err = generate(&buffer, pkg.Name, types)
					if err != nil {
						return fmt.Errorf("generate template: %w", err)
					}

					genFileName := fileName[:strings.Index(fileName, ".go")] + "_gen.go"
					src := buffer.Bytes()

					if b, err := format.Source(src); err != nil {
						_ = writeFile(src, genFileName)
						return fmt.Errorf("format template: %w", err)
					} else {
						return writeFile(b, genFileName)
					}
				}
			}
		}
	}

	return nil
}

func writeFile(b []byte, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	_ = file.Close()
	if err != nil {
		return fmt.Errorf("write template: %w", err)
	}
	return nil
}
