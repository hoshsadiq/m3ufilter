package m3u

import (
	"bufio"
	"github.com/PuerkitoBio/rehttp"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var log = logger.Get()

var client = NewClient(5)

func GetPlaylist(conf *config.Config) (streams Streams, allFailed bool) {
	streams = Streams{}

	errors := 0
	// todo we can do each provider in its own coroutine, then converged at the end.
	//   furthermore, each line can be done in its own coroutine as well.
	for _, provider := range conf.Providers {
		u, err := url.Parse(provider.Uri)
		if err != nil {
			errors++
			log.Errorf("Could not parse URL for %s, err = %v", provider.Uri, err)
			continue
		}

		if u.Scheme == "file" {
			log.Infof("reading from provider %s", u)
		} else {
			log.Infof("reading from provider %s://%s", u.Scheme, u.Host)
		}

		resp, err := client.Get(provider.Uri)
		if err != nil {
			errors++
			log.Errorf("could not retrieve playlist from provider %s, err = %v", provider.Uri, err)
			continue
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
			setOutputMarkers(conf.Core.Canonicalise.MainCountry, stream, nextStream)
		}
	}

	return streams, len(conf.Providers) == errors
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
		Timeout:   time.Second * 10,
		Transport: tr,
	}
}
