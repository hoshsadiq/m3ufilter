package server

import (
	"github.com/hoshsadiq/m3ufilter/m3u"
)

func updatePlaylist(conf *httpState) {
	if conf.appConfig.Core.AutoReloadConfig {
		err := conf.appConfig.Load()
		if err != nil {
			log.Errorf("Failed to reload the config, skipping updating playlist.")
			return
		}
	}

	conf.mut.Lock()
	defer conf.mut.Unlock()

	log.Info("Updating playlists")
	newPlaylists, newEpg, allFailed := m3u.ProcessConfig(conf.appConfig, false)
	if allFailed {
		log.Info("Skipping playlist synchronisation to server as all providers failed")
	} else {
		conf.playlists = &newPlaylists
		conf.epg = newEpg
	}
	log.Info("Done")
}
