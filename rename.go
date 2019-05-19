package main

import "github.com/grafov/m3u8"

func setSegmentValues(ms *m3u8.MediaSegment, setters []*Setter) {
	for _, setter := range setters {
		if shouldIncludeSegment(ms, setter.Filters) { // ensure this ANDed
			if setter.Name != "" {
				ms.Title = setter.Name
			}
			for attrKey, attrValue := range setter.Attributes {
				attr, err := getAttr(ms, attrKey)
				if err == nil {
					attr.Value = attrValue
				} else {
					ms.Attributes = append(ms.Attributes, &m3u8.Attribute{Key: attrKey, Value: attrValue})
				}
			}
		}
	}
}
