package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type FunctionInfo struct {
	Name    string
	Comment string
}

func ParseFunctionsFromDir(dir string) ([]FunctionInfo, error) {
	var funcs []FunctionInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" && !strings.HasSuffix(path, "_test.go") {
			fs := token.NewFileSet()
			node, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			for _, f := range node.Decls {
				if fn, ok := f.(*ast.FuncDecl); ok {
					funcs = append(funcs, FunctionInfo{
						Name:    fn.Name.Name,
						Comment: fn.Doc.Text(),
					})
				}
			}
		}
		return nil
	})

	return funcs, err
}
