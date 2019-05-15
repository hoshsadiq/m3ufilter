package main

import "regexp"

var reCache = make(map[string]*regexp.Regexp)

func getRegexCache(regex string) *regexp.Regexp {
	if reCache[regex] == nil {
		reCache[regex] = regexp.MustCompile(regex)
	}

	return reCache[regex]
}