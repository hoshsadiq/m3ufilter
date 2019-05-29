package m3u

func shouldIncludeSegment(segment *Stream, filters []string) bool {
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
