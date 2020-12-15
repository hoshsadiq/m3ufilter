package m3u

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/net"
	"net/http"
	"strings"
)

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
		resp, err = net.Head(stream.Uri)
	case "get":
		resp, err = net.Get(stream.Uri)
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
