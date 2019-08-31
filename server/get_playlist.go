package server

import (
	"errors"
	"github.com/hoshsadiq/m3ufilter/writer"
	"net/http"
)

func getPlaylist(conf *httpConfig, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		err := errors.New(http.StatusText(http.StatusMethodNotAllowed))
		return StatusError{Code: http.StatusMethodNotAllowed, Err: err}
	}

	w.Header().Set("Content-Type", "audio/mpegurl")

	writer.WriteOutput(conf.appConfig.Core.Output, w, *conf.playlists)
	return nil
}
