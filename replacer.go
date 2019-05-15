package main

import "github.com/grafov/m3u8"

type replaceAction struct {
	find string
	replace string
}

func doReplace(segment *m3u8.MediaSegment, replacements *Replacement) {
	if len(replacements.Name) > 0 {
		for _, replaceAction := range replacements.Name {
			var re = getRegexCache(replaceAction.Find)
			segment.Title = re.ReplaceAllString(segment.Title, replaceAction.Replace)
		}
	}

	if len(replacements.Attributes) == 0 {
		return
	}

	for attribKey, attrReplacements := range replacements.Attributes {
		attr, err := getAttr(segment, attribKey)
		if err != nil {
			continue
		}

		for _, replaceAction := range attrReplacements {
			var re = getRegexCache(replaceAction.Find)
			attr.Value = re.ReplaceAllString(attr.Value, replaceAction.Replace)
		}
	}
}

