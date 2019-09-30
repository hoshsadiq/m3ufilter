package m3u

import (
	"net/http"
)

// returns whether or not a stream should be included in the list
// A stream is considered include if any filter is matched. Once a filter matches the current stream, further filters
// will not be considered
// If no filters return true for the given stream, the stream is not included
// If a filter is matched, and the provider has been configured to be checked for failing links, it will not be returned
// if the stream does not survive a HEAD request or the content type isn't of a video type.
func shouldIncludeStream(stream *Stream, filters []string, checkStreams bool) bool {
	for _, filter := range filters {
		if filter == "" {
			continue
		}

		include, err := evaluateBool(stream, filter)
		if err != nil {
			log.Printf("error parsing expression %s, error = %v", filter, err)
		}

		if include {
			if checkStreams {
				return isWorkingStream(stream)
			}

			return true
		}
	}

	return len(filters) == 0
}

// Validates that a stream is indeed working.
// A stream is considered working if and only if the following are correct:
// We get a valid 200 response for a HEAD request to the stream endpoint
// The content-type of the request contains a valid stream
func isWorkingStream(stream *Stream) bool {
	resp, err := client.Head(stream.Uri)
	if err != nil {
		log.Errorf("Could not check stream is working for stream; id=%s, name=%s", stream.Id, stream.Name)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Stream did not return a valid response code; statusCode=%d, id=%s, name=%s", resp.StatusCode, stream.Id, stream.Name)
		return false
	}

	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "video/mp2t":
		fallthrough
	case "application/octet-stream":
		fallthrough
	case "application/vnd.apple.mpegurl":
		return true
	default:
		log.Errorf("Stream did not contain valid content type. If this a working stream, please report this as a bug; contentType=%s, id=%s, name=%s.", contentType, stream.Id, stream.Name)
		return false
	}
}
