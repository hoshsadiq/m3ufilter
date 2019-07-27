package m3u

import (
	"bufio"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/http"
	"github.com/hoshsadiq/m3ufilter/logger"
	"sort"
)

var log = logger.Get()

var client = http.NewClient(200, 5)

func GetPlaylist(conf *config.Config) Streams {
	streams := Streams{}
	// todo we can do each provider in its own coroutine, then converged at the end.
	//   furthermore, each line can be done in its own coroutine as well.
	for _, provider := range conf.Providers {
		log.Infof("reading from provider %s", provider.Uri)
		resp, err := client.Get(provider.Uri)
		if err != nil {
			log.Errorf("could not retrieve playlist from provider %s, err = %v", provider.Uri, err)
			continue
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Errorf("could not close request body for provider %s, err = %v", provider.Uri, err)
			}
		}()

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
