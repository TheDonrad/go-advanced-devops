package ststiclint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var osExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit()",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	// Анализируем только пакет main
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	osExitCheck := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if s, ok := call.Fun.(*ast.SelectorExpr); ok {
				if s.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "expression is os.Exit() in main func")
				}
			}
		}
	}

	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				return x.Name.Name == "main"
			case *ast.ExprStmt:
				osExitCheck(x)
			}
			return true
		})
	}
	return nil, nil
}
