package main

type IncludeExcludeFilter struct {
	Include []string
	Exclude []string
}

func (filter *IncludeExcludeFilter) shouldInclude(text string) bool {
	if len(filter.Include) == 0 {
		return !filter.shouldExclude(text)
	}

	for _, regex := range filter.Include {
		var re = getRegexCache(regex)
		if re.MatchString(text) {
			return true
		}
	}

	return !filter.shouldExclude(text)
}

func (filter *IncludeExcludeFilter) shouldExclude(text string) bool {
	for _, regex := range filter.Exclude {
		var re = getRegexCache(regex)
		if re.MatchString(text) {
			return true
		}
	}

	return len(filter.Exclude) == 0
}
