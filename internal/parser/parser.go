package router

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type APIInfo struct {
	Method       string
	Path         string
	AuthRequired bool
}

func ParseRouterFile(filepath string) (map[string]APIInfo, error) {
	apiMap := make(map[string]APIInfo)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filepath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if len(callExpr.Args) < 2 {
			return true
		}

		method := selExpr.Sel.Name // GET, POST, DELETE
		pathArg, _ := callExpr.Args[0].(*ast.BasicLit)
		handlerArg, _ := callExpr.Args[1].(*ast.SelectorExpr)

		if pathArg == nil || handlerArg == nil {
			return true
		}

		path := strings.Trim(pathArg.Value, "\"")
		handlerName := handlerArg.Sel.Name

		authRequired := false
		if len(callExpr.Args) > 2 {
			for _, extraArg := range callExpr.Args[2:] {
				if sel, ok := extraArg.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "Authenticate" {
						authRequired = true
					}
				}
			}
		}

		apiMap[handlerName] = APIInfo{
			Method:       method,
			Path:         path,
			AuthRequired: authRequired,
		}

		return true
	})

	return apiMap, nil
}
