package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"

	"github.com/fatih/structtag"
	"golang.org/x/tools/imports"

	"graphql-project/core"
)

type StructTypes = map[string]*ast.StructType

func FindStructTypes(file *ast.File) StructTypes {
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

func GetTagName(tags *structtag.Tags, key string) *string {
	if tags != nil {
		if tag, err := tags.Get(key); err == nil {
			return &tag.Name
		}
	}
	return nil
}

func GetTagOption(tags *structtag.Tags, key string, i int) *string {
	if tags != nil {
		if tag, err := tags.Get(key); err == nil {
			if i < len(tag.Options) {
				return &tag.Options[i]
			}
		}
	}
	return nil
}

func GetArg(i int) string {
	if i < len(os.Args) {
		return os.Args[i]
	}
	return ""
}

func Generate(pkgName string, pkgPath string, srcName string, generate func(string, string, StructTypes) error) error {

	file, err := parser.ParseFile(token.NewFileSet(), path.Join(pkgPath, srcName), nil, 0)
	if err != nil {
		return err
	}

	types := FindStructTypes(file)
	if len(types) > 0 {
		err = generate(srcName, pkgName, types)
		if err != nil {
			return err
		}
	}

	return nil
}

func Write(bytes []byte, fileName string) error {
	if b, err := imports.Process("", bytes, nil); err != nil {
		_ = writeFile(bytes, fileName)
		return fmt.Errorf("format template: %w", err)
	} else {
		return writeFile(b, fileName)
	}
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
