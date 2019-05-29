package cache

import "regexp"

var regexpCache = make(map[string]*regexp.Regexp)

func Regexp(regex string) *regexp.Regexp {
	if regexpCache[regex] == nil {
		regexpCache[regex] = regexp.MustCompile(regex)
	}

	return regexpCache[regex]
}
