package m3u

import "github.com/grafov/m3u8"

func ShouldIncludeSegment(segment *m3u8.MediaSegment, filters []string) bool {
	for _, filter := range filters {
		if filter == "" {
			continue
		}

		include, err := evaluateBool(segment, filter)
		if err != nil {
			panic(err) // todo actually we need to log this
		}

		if include {
			return true
		}
	}

	return false
}
