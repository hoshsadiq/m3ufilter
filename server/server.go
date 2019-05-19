package server

import (
	"bufio"
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"log"
	"net/http"
)

func Serve(conf *config.Config) {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	client := &http.Client{Transport: transport}

	http.HandleFunc("/playlist.m3u", func(w http.ResponseWriter, r *http.Request) {
		GetPlaylist(w, r, conf, client)
	})

	log.Fatal(http.ListenAndServe(conf.Core.Listen, nil))
}

func GetPlaylist(w http.ResponseWriter, r *http.Request, conf *config.Config, client *http.Client) {
	w.Header().Set("Content-Type", "audio/mpegurl")

	for _, provider := range conf.Providers {
		resp, err := client.Get(provider.Uri)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		p, listType, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), false)
		if err != nil {
			panic(err)
		}

		switch listType {
		case m3u8.MEDIA:
			mediapl := p.(*m3u8.MediaPlaylist)
			newp, err := m3u.ProcessPlaylist(mediapl, provider)
			if err != nil {
				panic(err)
			}
			_, err = w.Write(newp.Encode().Bytes())
			if err != nil {
				log.Println(err)
			}
		default:
			log.Printf("Found a non-media thing. Code needs to be updated!")
		}
	}
}
