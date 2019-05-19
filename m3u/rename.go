package m3u

import (
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/config"
)

func setSegmentValues(ms *m3u8.MediaSegment, setters []*config.Setter) {
	for _, setter := range setters {
		if ShouldIncludeSegment(ms, setter.Filters) { // ensure this ANDed
			if setter.Name != "" {
				ms.Title = setter.Name
			}
			for attrKey, attrValue := range setter.Attributes {
				attr, err := GetAttr(ms, attrKey)
				if err == nil {
					attr.Value = attrValue
				} else {
					ms.Attributes = append(ms.Attributes, &m3u8.Attribute{Key: attrKey, Value: attrValue})
				}
			}
		}
	}
}
