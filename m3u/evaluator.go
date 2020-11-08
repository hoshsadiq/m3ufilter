package m3u

import (
	"github.com/hoshsadiq/m3ufilter/cache"
)

func regexWordCallback(subject string, word string, callback func(string) string) string {
	re := cache.Regexp(`(?i)\b(` + word + `)\b`)

	subject = re.ReplaceAllStringFunc(subject, func(s string) string {
		return callback(s)
	})

	return subject
}

func removeWord(s string) string {
	return ""
}
