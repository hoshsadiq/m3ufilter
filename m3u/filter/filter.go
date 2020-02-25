package filter

import (
	"sort"
	"strings"
)

func EnsureUniqueUrls(urlsStr string) string {
	if ! strings.Contains(urlsStr, "|") {
		return urlsStr
	}

	urls := strings.Split(urlsStr, "|")
	sort.Strings(urls)
	j := 0
	for i := 1; i < len(urls); i++ {
		if urls[j] == urls[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		urls[j] = urls[i]
	}

	return strings.Join(urls[:j+1], "|")
}
