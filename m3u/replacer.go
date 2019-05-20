package m3u

import (
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/regex"
)

func replace(segment *m3u8.MediaSegment, replacements *config.Replacement) {
	if replacements == nil {
		return
	}

	if len(replacements.Name) > 0 {
		for _, replacer := range replacements.Name {
			var re = regex.GetCache(replacer.Find)
			newTitle := re.ReplaceAllString(segment.Title, replacer.Replace)
			if newTitle != segment.Title {
				log.Tracef("title %s replaced with %s; findReplace = %v", segment.Title, newTitle, replacer)
			}
			segment.Title = newTitle
			//attr := GetAttr(segment, "tvg-name")
			//attr.Value = segment.Title
		}
	}

	if len(replacements.Attributes) == 0 {
		return
	}

	for attribKey, attrReplacements := range replacements.Attributes {
		attr, err := GetAttr(segment, attribKey)
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
