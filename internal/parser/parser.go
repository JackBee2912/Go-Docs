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
	Name         string
	Comment      string
	RequestModel string
	ErrorCodes   []string
}

func ParseFunctionsFromDir(dir string) ([]FunctionInfo, error) {
	var funcs []FunctionInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		fileFuncs, err := parseFunctionsFromFile(path)
		if err != nil {
			return err
		}

		funcs = append(funcs, fileFuncs...)
		return nil
	})

	return funcs, err
}

func parseFunctionsFromFile(filepath string) ([]FunctionInfo, error) {
	var funcs []FunctionInfo

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil { // chỉ lấy method (có receiver, VD: func (h *Handler) XYZ)
			continue
		}

		funcInfo := FunctionInfo{
			Name:    fn.Name.Name,
			Comment: strings.TrimSpace(fn.Doc.Text()),
		}

		// Phân tích trong body
		if fn.Body != nil {
			for _, stmt := range fn.Body.List {
				inspectNode(stmt, &funcInfo)
			}
		}

		funcs = append(funcs, funcInfo)
	}

	return funcs, nil
}

func inspectNode(n ast.Node, info *FunctionInfo) {
	// Tìm request model: var req models.XYZ
	if declStmt, ok := n.(*ast.DeclStmt); ok {
		if genDecl, ok := declStmt.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if vspec, ok := spec.(*ast.ValueSpec); ok {
					for _, value := range vspec.Values {
						if callExpr, ok := value.(*ast.CallExpr); ok {
							if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
								if x, ok := fun.X.(*ast.Ident); ok && x.Name == "models" {
									info.RequestModel = "models." + fun.Sel.Name
								}
							}
						}
					}
				}
			}
		}
	}

	// Tìm error codes trong resp.BuildErrorResp(...)
	if exprStmt, ok := n.(*ast.ExprStmt); ok {
		if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
			if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "BuildErrorResp" && len(callExpr.Args) > 0 {
					if firstArg, ok := callExpr.Args[0].(*ast.SelectorExpr); ok {
						if x, ok := firstArg.X.(*ast.Ident); ok && x.Name == "resp" {
							errorCode := firstArg.Sel.Name
							info.ErrorCodes = append(info.ErrorCodes, errorCode)
						}
					}
				}
			}
		}
	}

	// Recursive check deeper blocks (if, for, etc.)
	ast.Inspect(n, func(child ast.Node) bool {
		switch c := child.(type) {
		case *ast.BlockStmt:
			for _, stmt := range c.List {
				inspectNode(stmt, info)
			}
		case *ast.IfStmt:
			inspectNode(c.Body, info)
			if c.Else != nil {
				inspectNode(c.Else, info)
			}
		case *ast.ForStmt:
			inspectNode(c.Body, info)
		}
		return true
	})
}
