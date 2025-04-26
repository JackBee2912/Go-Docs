package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// APIInfo chứa thông tin 1 API route
type APIInfo struct {
	Method       string // GET, POST, PUT, DELETE
	Path         string // /core/v1/friend/unfriend
	AuthRequired bool   // true nếu có middlewares.Authenticate
}

// ParseRouterFile phân tích file router.go và map tên handler function -> APIInfo
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

		method := selExpr.Sel.Name // GET, POST, PUT, DELETE...

		if len(callExpr.Args) < 2 {
			return true
		}

		pathLit, ok := callExpr.Args[0].(*ast.BasicLit)
		handlerSel, ok2 := callExpr.Args[1].(*ast.SelectorExpr)
		if !ok || !ok2 {
			return true
		}

		path := strings.Trim(pathLit.Value, "\"")
		handlerFunc := handlerSel.Sel.Name

		authRequired := false
		if len(callExpr.Args) >= 3 {
			for _, arg := range callExpr.Args[2:] {
				if sel, ok := arg.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "Authenticate" {
						authRequired = true
					}
				}
			}
		}

		apiMap[handlerFunc] = APIInfo{
			Method:       method,
			Path:         path,
			AuthRequired: authRequired,
		}

		return true
	})

	return apiMap, nil
}
