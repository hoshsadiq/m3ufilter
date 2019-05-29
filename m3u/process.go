package m3u

import (
	"bufio"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"net/http"
)

var log = logger.Get()

func GetPlaylist(conf *config.Config) []*Stream {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	client := &http.Client{Transport: transport}

	playlists := []*Stream{}
	for _, provider := range conf.Providers {
		log.Infof("reading from provder %s", provider.Uri)
		resp, err := client.Get(provider.Uri)
		if err != nil {
			log.Errorf("could not retrieve playlist from provider %s, err = %v", provider.Uri, err)
		}
		defer resp.Body.Close()

		pl, err := decode(bufio.NewReader(resp.Body), provider)
		if err != nil {
			log.Errorf("could not retrieve playlist from provider %s, err = %v", provider.Uri, err)
		}

		playlists = append(playlists, pl...)
	}

	return playlists
}
