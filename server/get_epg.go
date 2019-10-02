package server

import (
	"errors"
	"net/http"
)

func getEpg(conf *httpConfig, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		err := errors.New(http.StatusText(http.StatusMethodNotAllowed))
		return StatusError{Code: http.StatusMethodNotAllowed, Err: err}
	}

	// todo or is it: text/xml
	w.Header().Set("Content-Type", "application/xml")

	//writer.WriteOutput(conf.appConfig.Core.Output, w, *conf.Epg)
	return nil
}
