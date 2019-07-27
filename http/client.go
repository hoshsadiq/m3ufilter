package http

import (
	"net/http"
	"time"
)

// Naive Retry - every 2 seconds
func NewClient(RequiredResponseCode int, MaxRetryAttempts int) *http.Client {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	return &http.Client{
		Timeout: time.Second * 3,
		Transport: &Retry{
			RoundTripper:        transport,
			RequireResponseCode: RequiredResponseCode,
			MaxAttempts:         MaxRetryAttempts,
		},
	}
}
