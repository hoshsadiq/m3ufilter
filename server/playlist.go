package server

import (
	"github.com/hoshsadiq/m3ufilter/m3u"
)

func updatePlaylist(conf *httpState) {
	if conf.lock {
		log.Info("Retrieval is locked, trying again next time...")
		return
	}

	if conf.appConfig.Core.AutoReloadConfig {
		err := conf.appConfig.Load()
		if err != nil {
			log.Errorf("Failed to reload the config, skipping updating playlist.")
			return
		}
	}

	conf.lock = true
	defer func() {
		conf.lock = false
	}()

	log.Info("Updating playlists")
	newPlaylists, newEpg, allFailed := m3u.ProcessConfig(conf.appConfig)
	if allFailed {
		log.Info("Skipping playlist synchronisation to server as all providers failed")
	} else {
		conf.playlists = &newPlaylists
		conf.epg = newEpg
	}
	log.Info("Done")
}
