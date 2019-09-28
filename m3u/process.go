package m3u

import (
	"bufio"
	"github.com/PuerkitoBio/rehttp"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"net/http"
	"net/url"
	"sort"
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

		log.Infof("reading from provider %s://%s", u.Scheme, u.Host)
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

	return streams, len(conf.Providers) == errors
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
