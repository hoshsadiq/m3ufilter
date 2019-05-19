package regex

import "regexp"

var cache = make(map[string]*regexp.Regexp)

func GetCache(regex string) *regexp.Regexp {
	if cache[regex] == nil {
		cache[regex] = regexp.MustCompile(regex)
	}

	return cache[regex]
}
