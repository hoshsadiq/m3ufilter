package m3u

import (
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/util"
)

func setSegmentValues(ms *m3u8.MediaSegment, setters []*config.Setter, syncTitleName bool) {
	for _, setter := range setters {
		if shouldIncludeSegment(ms, setter.Filters) {
			if setter.Name != "" {
				newTitle, err := evaluateStr(ms, setter.Name)
				if err != nil {
					log.Errorln(err)
				}
				if newTitle != ms.Title {
					log.Tracef("title %s replaced with %s; expr = %v", ms.Title, newTitle, setter.Name)
				}

				ms.Title = newTitle
				if syncTitleName {
					util.SetAttr(ms, "tvg-name", newTitle)
				}
			}
			for attrKey, attrValue := range setter.Attributes {
				newValue, err := evaluateStr(ms, attrValue)
				if err != nil {
					log.Errorln(err)
				}

				util.SetAttr(ms, attrKey, newValue)
			}
		}
	}
}
