package m3u

import (
	"bufio"
	"github.com/PuerkitoBio/rehttp"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u/xmltv"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var log = logger.Get()

var client = NewClient(5)

func GetPlaylist(conf *config.Config) (streams Streams, allFailed bool) {
	errors := 0

	epg, err := getEpg(conf.EpgProviders)
	if err != nil {
		panic(err) // todo
	}
	//log.Printf("%v", epg)

	streams = Streams{}
	// todo we can do each provider in its own coroutine, then converged at the end.
	//   furthermore, each line can be done in its own coroutine as well.
	for _, provider := range conf.Providers {
		resp, err := getUri(provider.Uri)
		if err != nil {
			errors++
			log.Errorf("Could not retrieve file from %s, err = %v", provider.Uri, err)
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Errorf("could not close request body for provider %s, err = %v", provider.Uri, err)
			}
		}()

		pl, err := decode(conf, bufio.NewReader(resp.Body), provider)
		if err != nil {
			errors++
			log.Errorf("could not decode playlist from provider %s, err = %v", provider.Uri, err)
			continue
		} else {
			streams = append(streams, pl...)
		}
	}

	sort.Sort(streams)

	if conf.Core.Canonicalise.Enable {
		streamsLength := len(streams)
		var nextStream *Stream
		for i, stream := range streams {
			if i+1 >= streamsLength {
				continue
			}

			nextStream = streams[i+1]
			setOutputMarkers(conf.Core.Canonicalise.DefaultCountry, stream, nextStream)
		}
	}

	for _, stream := range streams {
		for _, epgChannel := range epg.Channel {
			if stream.meta.originalId == epgChannel.ID {
				epgChannel.DisplayName = append(epgChannel.DisplayName, xmltv.DisplayName{Value:stream.GetName()})
				break
			}
			//if stream.meta.originalName == epgChannel.DisplayName {
			//
			//}
			log.Tracef("%v,%v", stream, epg)
		}
	}

	return streams, len(conf.Providers) == errors
}

func getUri(uri string) (*http.Response, error) {
	u, err := url.Parse(uri)
	if err != nil {
		log.Errorf("Could not parse URL for %s, err = %v", uri, err)
		return nil, err
	}
	if u.Scheme == "file" {
		log.Infof("reading from %s", u)
	} else {
		log.Infof("reading from %s://%s", u.Scheme, u.Host)
	}
	resp, err := client.Get(uri)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func getEpg(providers []*config.EpgProvider) (*xmltv.XMLTV, error) {
	var epg xmltv.XMLTV

	for _, provider := range providers {
		resp, err := getUri(provider.Uri)
		if err != nil {
			return nil, err
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Errorf("could not close request body for provider %s, err = %v", provider.Uri, err)
			}
		}()

		err = xmltv.Load(resp.Body, &epg)
		if err != nil {
			return nil, err
		}
	}

	var channels = make(map[string]*xmltv.Channel, len(epg.Channel))
	var nameIdMapping = make(map[string]string)

	for _, c := range epg.Channel {
		channel, ok := channels[c.ID]
		var found = false
		if !ok {
			for _, dpname := range c.DisplayName {
				channelId, ok := nameIdMapping[dpname.Value]
				if ok {
					channel, ok = channels[channelId]
					nameIdMapping[c.ID] = channelId
					found = true
					break
				}
			}

			if !found {
				channels[c.ID] = c
				for _, dpname := range c.DisplayName {
					nameIdMapping[dpname.Value] = c.ID
				}
				continue
			}
		}
		for _, left := range c.DisplayName {
			found = false
			for _, right := range channel.DisplayName {
				if right.Value == left.Value {
					found = true
					break
				}
			}
			if !found {
				channel.DisplayName = append(channel.DisplayName, left)
			}

			nameIdMapping[left.Value] = c.ID
		}
	}

	return &epg, nil
}

func setOutputMarkers(mainCountry string, left *Stream, right *Stream) {
	if left.meta.canonicalName != right.meta.canonicalName {
		return
	}

	if left.meta.country != right.meta.country {
		mainCountry = strings.ToUpper(mainCountry)
		if left.meta.country == "" || left.meta.country != mainCountry {
			left.meta.showCountry = true
		}

		if right.meta.country == "" || right.meta.country != mainCountry {
			right.meta.showCountry = true
		}
	} else {
		left.meta.showDefinition = left.meta.definition != right.meta.definition
		right.meta.showDefinition = left.meta.showDefinition
	}
}

func NewClient(MaxRetryAttempts int) *http.Client {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	tr := rehttp.NewTransport(
		transport,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(MaxRetryAttempts),
			rehttp.RetryStatuses(200),
			rehttp.RetryTemporaryErr(),
		),
		rehttp.ConstDelay(time.Second),
	)
	return &http.Client{
		Timeout:   time.Second * 30,
		Transport: tr,
	}
}
