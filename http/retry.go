package http

import (
	"errors"
	"github.com/hoshsadiq/m3ufilter/logger"
	"net/http"
	"strconv"
	"time"
)

var log = logger.Get()

// Retry - http client with retry support
type Retry struct {
	http.RoundTripper

	MaxAttempts         int
	RequireResponseCode int
}

// Naive Retry - every 2 seconds
func (r *Retry) RoundTrip(req *http.Request) (*http.Response, error) {
	attempts := 0

	if r.MaxAttempts < 1 {
		return nil, errors.New("MaxAttempts must be at least 1")
	}

	for {
		attempts++
		resp, err := r.RoundTripper.RoundTrip(req)

		if err == nil && ((r.RequireResponseCode == 0 && resp.StatusCode < 500) || (r.RequireResponseCode == resp.StatusCode)) {
			return resp, err
		}

		log.Infof("Request attempt %d failed, err = %v", attempts, err)

		delay, retry := r.shouldRetry(attempts, resp)
		if !retry {
			log.Infof("Request failed, not retrying again.")
			return resp, err
		}

		log.Infof("Retrying in %s", delay)
		select {
		// check if canceled or timed-out
		case <-req.Context().Done():
			return resp, req.Context().Err()
		case <-time.After(delay):
		}
	}
}

func (r *Retry) shouldRetry(attempts int, response *http.Response) (time.Duration, bool) {
	if attempts >= r.MaxAttempts {
		return time.Duration(0), false
	}

	delay := time.Duration(attempts) * time.Second

	if response != nil && response.Header.Get("Retry-After") != "" {
		after, err := strconv.ParseInt(response.Header.Get("Retry-After"), 10, 64)
		if err != nil && after > 0 {
			delay = time.Duration(after) * time.Second
		}
	}

	return delay, true
}
