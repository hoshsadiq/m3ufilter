package cache

import (
	"github.com/hoshsadiq/govaluate"
)

var exprCache = make(map[string]*govaluate.EvaluableExpression)

func Expr(expr string, functions map[string]govaluate.ExpressionFunction) (*govaluate.EvaluableExpression, error) {
	if exprCache[expr] == nil {
		evaluableExpression, err := govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
		if err != nil {
			return nil, err
		}
		exprCache[expr] = evaluableExpression
	}

	return exprCache[expr], nil
}
