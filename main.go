package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	v := visitor{fset: token.NewFileSet()}
	arr := []string{"./test/test_1.go"}

	for _, filePath := range arr {
		if filePath == "--" { // to be able to run this like "go run main.go -- input.go"
			continue
		}

		f, err := parser.ParseFile(v.fset, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}
		ast.Walk(&v, f)
	}
}

type visitor struct {
	fset *token.FileSet
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}

	println("func check:", funcDecl.Name.Name)
	for _, field := range funcDecl.Type.Params.List {
		var fieldSelectorExpr *ast.SelectorExpr
		logPointer := false
		fieldTypePointer, ok := field.Type.(*ast.StarExpr)
		if !ok {
			fieldSelectorExpr, ok = field.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}
		} else {
			fieldSelectorExpr, ok = fieldTypePointer.X.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			logPointer = true
		}

		fieldIdentPackage, ok := fieldSelectorExpr.X.(*ast.Ident)
		if !ok {
			continue
		}

		println("logger name:", field.Names[0].Name)
		println("logger package:", fieldIdentPackage.Name)
		println("logger name:", fieldSelectorExpr.Sel.Name)
		println("is pointer:", logPointer)
		println()
	}

	//var buf bytes.Buffer
	//printer.Fprint(&buf, v.fset, node)
	////println()
	//fmt.Printf("%s | %#v\n", buf.String(), node)
	return v
}
