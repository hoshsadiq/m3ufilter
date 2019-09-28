package server

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"github.com/mileusna/crontab"
	"net/http"
)

var log = logger.Get()

func Serve(appConfig *config.Config) {
	conf := &httpConfig{
		playlists: &m3u.Streams{},
		lock:      false,
		appConfig: appConfig,
		crontab:   crontab.New(),
	}

	log.Info("Scheduling cronjob to periodically update playlist.")
	scheduleJob(conf, appConfig.Core.UpdateSchedule)

	log.Info("Parsing for the first time...")
	conf.crontab.RunAll()

	log.Info("Starting server")
	http.Handle("/playlist.m3u", httpHandler{conf, getPlaylist})
	http.Handle("/epg.xml", httpHandler{conf, getEpg})
	http.Handle("/update", httpHandler{conf, postUpdate})

	server := &http.Server{Addr: appConfig.Core.ServerListen}
	log.Fatal(server.ListenAndServe())
}

func scheduleJob(conf *httpConfig, schedule string) {
	conf.crontab.MustAddJob(schedule, func() {
		updatePlaylist(conf)
	})
}
