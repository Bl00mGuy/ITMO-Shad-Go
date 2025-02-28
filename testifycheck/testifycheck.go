//go:build !solution

package testifycheck

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "require",
	Doc:  "Analyzes code for improper error handling using testify",
	Run:  analyzeTestifyUsage,
}

func isErrorType(expression ast.Expr, pass *analysis.Pass) bool {
	expressionType := pass.TypesInfo.TypeOf(expression)
	interfaceType, isInterface := expressionType.Underlying().(*types.Interface)
	if !isInterface {
		return false
	}

	return interfaceType.NumMethods() == 1 && interfaceType.Method(0).FullName() == "(error).Error"
}

func analyzeTestifyUsage(pass *analysis.Pass) (interface{}, error) {
	errorHandlingAlternatives := map[string]string{
		"Nil":     "NoError",
		"Nilf":    "NoErrorf",
		"NotNil":  "Error",
		"NotNilf": "Errorf",
	}

	for fileIndex, sourceFile := range pass.Files {
		if fileIndex == 1 {
			continue
		}

		ast.Inspect(sourceFile, func(node ast.Node) bool {
			callExpression, isCallExpr := node.(*ast.CallExpr)
			if !isCallExpr {
				return true
			}

			calledFunction, _ := typeutil.Callee(pass.TypesInfo, callExpression).(*types.Func)
			if calledFunction == nil {
				return true
			}

			packageName := calledFunction.Pkg().Name()
			arguments := callExpression.Args
			if len(arguments) < 1 ||
				(!isErrorType(arguments[0], pass) && (len(arguments) < 2 || !isErrorType(arguments[1], pass))) {
				return true
			}

			if alternativeFunctionName, exists := errorHandlingAlternatives[calledFunction.Name()]; exists {
				pass.Reportf(
					callExpression.Pos(),
					"use %s.%s instead of comparing error to nil",
					packageName,
					alternativeFunctionName,
				)
				return false
			}
			return true
		})
	}
	return nil, nil
}
