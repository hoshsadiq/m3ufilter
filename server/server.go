package server

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"net/http"
)

var log = logger.Get()

func Serve(conf *config.Config) {
	log.Info("starting server")

	http.HandleFunc("/playlist.m3u", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/mpegurl")

		m3u.GetPlaylist(w, conf)
	})

	server := &http.Server{Addr: conf.Core.ServerListen}
	log.Fatal(server.ListenAndServe())
}
