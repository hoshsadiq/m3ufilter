package m3u

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"net/http"
	"strings"
)

// returns whether or not a stream should be included in the list
// A stream is considered include if any filter is matched. Once a filter matches the current stream, further filters
// will not be considered
// If no filters return true for the given stream, the stream is not included
// If a filter is matched, and the provider has been configured to be checked for failing links, it will not be returned
// if the stream does not survive a HEAD request or the content type isn't of a video type.
func shouldIncludeStream(stream *Stream, filters []string, checkStreams config.CheckStreams) bool {
	isWorkingStream := isWorkingStream(stream, checkStreams)
	for _, filter := range filters {
		if filter == "" {
			continue
		}

		include, err := evaluateBool(stream, filter)
		if err != nil {
			log.Printf("error parsing expression %s, error = %v", filter, err)
		}

		if include {
			return isWorkingStream
		}
	}

	return len(filters) == 0 && isWorkingStream
}

// Validates that a stream is indeed working.
// A stream is considered working if and only if the following are correct:
// We get a valid 200 response for a HEAD request to the stream endpoint
// The content-type of the request contains a valid stream
func isWorkingStream(stream *Stream, checkSteamConfig config.CheckStreams) bool {
	if stream.meta.available || !checkSteamConfig.Enabled {
		return true
	}

	stream.meta.available = true

	var resp *http.Response
	var err error
	switch strings.ToLower(checkSteamConfig.Method) {
	case "head":
		resp, err = client.Head(stream.Uri)
	case "get":
		resp, err = client.Get(stream.Uri)
	default:
		log.Errorf("provider.check_streams.method can only be head or get, got: %s", checkSteamConfig.Method)
		return true
	}

	if err != nil {
		log.Errorf("Could not check stream is working for stream; id=%s, name=%s", stream.Id, stream.Name)
		stream.meta.available = false
		return checkSteamConfig.Action != config.InvalidStreamRemove
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Stream did not return a valid response code; statusCode=%d, id=%s, name=%s", resp.StatusCode, stream.Id, stream.Name)
		stream.meta.available = false
		return checkSteamConfig.Action != config.InvalidStreamRemove
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
		stream.meta.available = false
		return checkSteamConfig.Action != config.InvalidStreamRemove
	}
}
