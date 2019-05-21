package m3u

import (
	"bufio"
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"io"
	"net/http"
)

var log = logger.Get()

func GetPlaylist(w io.Writer, conf *config.Config) {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	client := &http.Client{Transport: transport}

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
			newp, err := processPlaylist(mediapl, provider, conf.Core.SyncTitleName)
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

func processPlaylist(pl *m3u8.MediaPlaylist, providerConfig *config.Provider, syncTitleName bool) (*m3u8.MediaPlaylist, error) {
	p, err := m3u8.NewMediaPlaylist(pl.Count(), pl.Count())
	if err != nil {
		return nil, err
	}

	for _, segment := range pl.Segments {
		if segment == nil {
			continue
		}

		segment = segment.Clone()
		//segment.Title, err = strconv.Unquote("\"" + segment.Title + "\"")
		//if err != nil {
		//	log.Errorf("error unquoting %s", segment.Title)
		//}

		replace(segment, providerConfig.Replacements, syncTitleName)

		if !shouldIncludeSegment(segment, providerConfig.Filters) {
			continue
		}

		setSegmentValues(segment, providerConfig.Setters, syncTitleName)
		//renameGroups(segment, providerConfig.Groups)

		err = p.AppendSegment(segment)
		if err != nil {
			return nil, err
		}
	}

	err = p.SetWinSize(p.Count())
	if err != nil {
		return nil, err
	}

	return p, nil
}
