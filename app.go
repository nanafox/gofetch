package gofetch

import (
	"net/http"
	"time"
)

// New creates and returns a new instance of Client with
// optional timeout and debug settings
func New(config ...Config) (client *Client) {
	if config != nil {
		client = &Client{Config: Config{Timeout: config[0].Timeout, Debug: config[0].Debug}}
	} else {
		client = &Client{Config: Config{}}
	}

	if client.Config.Timeout == 0 {
		// 500 ms Default timeout if not provided
		client.Config.Timeout = 500 * time.Millisecond
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,              // Number of idle connections to maintain
			MaxIdleConnsPerHost: 10,               // Max number of idle connections per host
			IdleConnTimeout:     90 * time.Second, // Timeout for idle connections
			DisableKeepAlives:   false,            // Keep connections alive
		},
		Timeout: client.Config.Timeout,
	}

	client.httpClient = httpClient
	return
}
