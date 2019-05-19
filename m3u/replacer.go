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
		for _, replaceAction := range replacements.Name {
			var re = regex.GetCache(replaceAction.Find)
			segment.Title = re.ReplaceAllString(segment.Title, replaceAction.Replace)
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

		for _, replaceAction := range attrReplacements {
			var re = regex.GetCache(replaceAction.Find)
			attr.Value = re.ReplaceAllString(attr.Value, replaceAction.Replace)
		}
	}
}
