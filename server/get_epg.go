package server

import (
	"errors"
	"github.com/hoshsadiq/m3ufilter/m3u/xmltv"
	"net/http"
)

func getEpg(state *httpState, w http.ResponseWriter, r *http.Request) error {
	var err error
	if r.Method != "GET" {
		err = errors.New(http.StatusText(http.StatusMethodNotAllowed))
		return StatusError{Code: http.StatusMethodNotAllowed, Err: err}
	}

	w.Header().Set("Content-Type", "application/xml")
	return xmltv.Dump(w, state.epg, state.appConfig.Core.PrettyOutputXml)
}
