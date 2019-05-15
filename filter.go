package main

import "github.com/grafov/m3u8"

func shouldIncludeSegment(segment *m3u8.MediaSegment, filters *Filters) bool {
	if filters.Name != nil {
		if !filters.Name.shouldInclude(segment.Title) {
			return false
		}
	}

	if len(filters.Attributes) == 0 {
		return true
	}

	for attrib, filter := range filters.Attributes {
		attrib, _ := getAttr(segment, attrib)
		if !filter.shouldInclude(attrib.Value) {
			return false
		}
	}

	return true
}
