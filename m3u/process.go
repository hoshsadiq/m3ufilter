package m3u

import (
	"bufio"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"net/http"
	"sort"
	"time"
)

var log = logger.Get()

var client *http.Client

func GetPlaylist(conf *config.Config) Streams {
	createClientIfNotExists()

	streams := Streams{}
	for _, provider := range conf.Providers {
		log.Infof("reading from provider %s", provider.Uri)
		resp, err := client.Get(provider.Uri)
		if err != nil {
			log.Errorf("could not retrieve playlist from provider %s, err = %v", provider.Uri, err)
			continue
		}
		defer resp.Body.Close()

		pl, err := decode(bufio.NewReader(resp.Body), provider)
		if err != nil {
			log.Errorf("could not decode playlist from provider %s, err = %v", provider.Uri, err)
			continue
		}

		streams = append(streams, pl...)
	}

	sort.Sort(streams)

	return streams
}

func createClientIfNotExists() {
	if client == nil {
		transport := &http.Transport{}
		transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		client = &http.Client{
			Timeout:   time.Second * 3,
			Transport: transport,
		}
	}
}
