// Package osexit анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main.
package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OSexitCheckAnalyzer переменная анализатора прямой вызов os.Exit в функции main пакета main.
var OSexitCheckAnalyzer = &analysis.Analyzer{
	Name: "osexitCheck",
	Doc:  "check for use os.Exit in main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			// функцией ast.Inspect проходим по всем узлам AST
			ast.Inspect(file, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.FuncDecl:
					if x.Name.Name != "main" {
						return false
					}
				case *ast.SelectorExpr:
					if x.Sel.Name == "Exit" {
						pass.Reportf(node.Pos(), "call os.Exit within main package and function")
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
