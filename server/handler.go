package server

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"github.com/hoshsadiq/m3ufilter/m3u/xmltv"
	"github.com/mileusna/crontab"
	"net/http"
)

type httpState struct {
	appConfig *config.Config
	playlists *m3u.Streams
	lock      bool
	crontab   *crontab.Crontab
	epg       *xmltv.XMLTV
}

type httpHandler struct {
	conf     *httpState
	callback func(e *httpState, w http.ResponseWriter, r *http.Request) error
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.callback(h.conf, w, r)
	if err != nil {
		switch e := err.(type) {
		case HttpError:
			log.Infof("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			log.Errorf("HTTP error retrieved %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}
