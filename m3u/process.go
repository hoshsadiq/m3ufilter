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

		pl, err := decode(bufio.NewReader(resp.Body))

		newp, err := processPlaylist(pl, provider)
		if err != nil {
			log.Errorf("unable to parse %s, err = %v", provider.Uri, err)
			continue
		}

		playlists = append(playlists, newp...)
	}

	return playlists
}

func processPlaylist(streams []*Stream, providerConfig *config.Provider) ([]*Stream, error) {
	newStreams := []*Stream{}

	for _, ms := range streams {
		if ms == nil {
			continue
		}
		log.Debugf("Processing segment: tvg-id=%s; channel=%s", ms.Id, ms.Name)

		if !shouldIncludeSegment(ms, providerConfig.Filters) {
			continue
		}

		setSegmentValues(ms, providerConfig.Setters)

		newStreams = append(newStreams, ms)
	}

	return newStreams, nil
}
