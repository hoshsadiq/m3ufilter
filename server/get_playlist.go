package server

import (
	"errors"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/writer"
	"net/http"
)

func getPlaylist(state *httpState, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "HEAD" && r.Method != "GET" {
		logger.Get().Errorf("Method %s is not allowed", r.Method)
		err := errors.New(http.StatusText(http.StatusMethodNotAllowed))
		return StatusError{Code: http.StatusMethodNotAllowed, Err: err}
	}

	w.Header().Set("Content-Type", "audio/mpegurl")

	writer.WriteOutput(state.appConfig.Core.Output, w, *state.playlists)
	return nil
}
