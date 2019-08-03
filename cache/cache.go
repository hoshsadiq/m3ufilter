package cache

import (
	"github.com/maja42/no-comment"
	"regexp"
	"sync"
)

var cache sync.Map

func Expr(expr string) string {
	key := "expr:" + expr
	if x, found := cache.Load(key); found {
		return x.(string)
	}

	uncommented := nocomment.StripCStyleComments(expr)
	cache.Store(key, uncommented)

	return uncommented
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
