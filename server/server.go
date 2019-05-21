package server

import (
	"bufio"
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"net/http"
)

var log = logger.Get()

func Serve(conf *config.Config) {
	log.Info("starting server")

	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	client := &http.Client{Transport: transport}

	http.HandleFunc("/playlist.m3u", func(w http.ResponseWriter, r *http.Request) {
		GetPlaylist(w, r, conf, client)
	})

	server := &http.Server{Addr: conf.Core.Listen}
	log.Fatal(server.ListenAndServe())
}

func GetPlaylist(w http.ResponseWriter, r *http.Request, conf *config.Config, client *http.Client) {
	w.Header().Set("Content-Type", "audio/mpegurl")

	for _, provider := range conf.Providers {
		log.Infof("reading from provder %s", provider.Uri)
		resp, err := client.Get(provider.Uri)
		if err != nil {
			log.Errorf("could not retrieve playlist from provider %s, err = %v", provider.Uri, err)
		}
		defer resp.Body.Close()

		p, listType, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), false)
		if err != nil {
			log.Errorf("could not parse provider %s, err = %v", provider.Uri, err)
			continue
		}

		switch listType {
		case m3u8.MEDIA:
			log.Debugf("found media type %v", listType)
			mediapl := p.(*m3u8.MediaPlaylist)
			newp, err := m3u.ProcessPlaylist(mediapl, provider, conf.Core.SyncTitleName)
			if err != nil {
				log.Errorf("unable to parse %s, err = %v", provider.Uri, err)
				continue
			}
			_, err = w.Write(newp.Encode().Bytes())
			if err != nil {
				log.Errorf("unable to write new playlist, err = %v", provider.Uri, err)
			}
		default:
			log.Errorf("found unsupported media type. code needs to be updated. Type = %v, err = %v", listType, err)
		}
	}
}
