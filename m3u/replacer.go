package m3u

import (
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/regex"
)

func replace(ms *m3u8.MediaSegment, replacements *config.Replacement, syncTitleName bool) {
	if replacements == nil {
		return
	}

	if len(replacements.Name) > 0 {
		for _, replacer := range replacements.Name {
			var re = regex.GetCache(replacer.Find)
			newTitle := re.ReplaceAllString(ms.Title, replacer.Replace)
			if newTitle != ms.Title {
				log.Tracef("title %s replaced with %s; findReplace = %v", ms.Title, newTitle, replacer)
			}
			ms.Title = newTitle
			if syncTitleName {
				SetAttr(ms, "tvg-name", newTitle)
			}
		}
	}

	if len(replacements.Attributes) == 0 {
		return
	}

	for attribKey, attrReplacements := range replacements.Attributes {
		attr, err := GetAttr(ms, attribKey)
		if err != nil {
			continue
		}

		for _, replacer := range attrReplacements {
			var re = regex.GetCache(replacer.Find)
			newValue := re.ReplaceAllString(attr.Value, replacer.Replace)
			if newValue != attr.Value {
				log.Tracef("attr %v has new value = %s; findReplace = %v", attr, newValue, replacer)
			}
		}
	}
}
