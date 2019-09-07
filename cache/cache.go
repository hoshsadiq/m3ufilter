package cache

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"regexp"
	"sync"
)

var cache sync.Map

func Expr(expression string, envArg interface{}) (*vm.Program, error) {
	key := "expr:" + expression
	if x, found := cache.Load(key); found {
		return x.(*vm.Program), nil
	}


	program, err := expr.Compile(expression, expr.Env(envArg))
	if err != nil {
		return nil, err
	}

	cache.Store(key, program)

	return program, nil
}

func Regexp(regex string) *regexp.Regexp {
	key := "regexp:" + regex

	if x, found := cache.Load(key); found {
		return x.(*regexp.Regexp)
	}
	re := regexp.MustCompile(regex)
	cache.Store(key, re)

	return re
}
