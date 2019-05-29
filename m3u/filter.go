package m3u

import (
	"github.com/grafov/m3u8"
)

func shouldIncludeSegment(segment *m3u8.MediaSegment, filters []string) bool {
	for _, filter := range filters {
		if filter == "" {
			continue
		}

		include, err := evaluateBool(segment, filter)
		if err != nil {
			log.Printf("error parsing expression %s, error = %v", filter, err)
		}

		if include {
			return true
		}
	}

	return false
}
