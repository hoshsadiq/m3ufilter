package server

import (
	"errors"
	"net/http"
)

func postUpdate(state *httpState, _ http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		err := errors.New(http.StatusText(http.StatusMethodNotAllowed))
		return StatusError{Code: http.StatusMethodNotAllowed, Err: err}
	}

	state.crontab.RunAll()

	err := errors.New(http.StatusText(http.StatusNoContent))
	return StatusError{Code: http.StatusNoContent, Err: err}
}
