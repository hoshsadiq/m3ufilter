package cache

import (
	nocomment "github.com/maja42/no-comment"
)

var exprCache = make(map[string]*string)

func Expr(expr string) string {
	if exprCache[expr] == nil {
		uncommented := nocomment.StripCStyleComments(expr)
		exprCache[expr] = &uncommented
	}

	return *exprCache[expr]
}
