package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "uniquelog",
	Doc:      "Checks that all logs and errors have unique message",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var logs = make(map[string]LogOutput)

type LogOutput struct {
	Position token.Pos
	File     *token.File
}

func run(pass *analysis.Pass) (interface{}, error) {
	// TODO check constant handling
	// TODO add to vars, add prebuild loggers, add different loggers support
	loggerPackage := "slog"
	loggerName := "Logger"

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	insp.Preorder(nodeFilter, func(node ast.Node) {
		var funcLogger string
		// parse func definition
		if node == nil {
			return
		}

		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return
		}

		println("func check:", funcDecl.Name.Name)

		// try to find logger in func args
		loggerFounded := false
		for _, field := range funcDecl.Type.Params.List {
			var fieldSelectorExpr *ast.SelectorExpr
			//logPointer := false
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
				//logPointer = true
			}

			fieldIdentPackage, ok := fieldSelectorExpr.X.(*ast.Ident)
			if !ok {
				continue
			}

			if fieldIdentPackage.Name == loggerPackage && fieldSelectorExpr.Sel.Name == loggerName {
				funcLogger = field.Names[0].Name
				loggerFounded = true
			}
			//println("logger name:", field.Names[0].Name)
			//println("logger package:", fieldIdentPackage.Name)
			//println("logger type:", fieldSelectorExpr.Sel.Name)
			//println("is pointer:", logPointer)
			//println()
		}
		if !loggerFounded {
			return
		}
		println("func logger:", funcLogger)

		for _, stmt := range funcDecl.Body.List {
			expression, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}
			callExpr, ok := expression.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			callExprFun, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			callExprFunIdent, ok := callExprFun.X.(*ast.Ident)
			if !ok {
				continue
			}
			if funcLogger != callExprFunIdent.Name {
				continue
			}
			for _, arg := range callExpr.Args {
				argBasicLit, ok := arg.(*ast.BasicLit)
				if !ok {
					continue
				}
				if foundedLog, ok := logs[argBasicLit.Value]; ok {
					pass.Reportf(argBasicLit.Pos(), "dublicated log: %v %v", argBasicLit.Value, foundedLog.FileLocation())
				} else {
					logs[argBasicLit.Value] = LogOutput{
						Position: argBasicLit.Pos(),
						File:     pass.Fset.File(argBasicLit.Pos()),
					}
				}
			}
		}
	})
	return nil, nil
}

func (l LogOutput) FileLocation() string {
	return fmt.Sprintf("%v:%v", l.File.Name(), l.File.Line(l.Position))
}
