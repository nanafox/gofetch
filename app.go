package gofetch

import (
	"net/http"
	"time"
)

// New creates and returns a new instance of Client with
// optional timeout and debug settings
func New(cfg Config) *Client {
	var c *Client

	if cfg.Timeout == 0 {
		cfg.Timeout = 500 * time.Millisecond
	}
	c = &Client{Config: cfg}

	if c.Config.Timeout == 0 {
		// 500ms Default timeout if not provided
		c.Config.Timeout = 500 * time.Millisecond
	}

	hc := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,              // Number of idle connections to maintain
			MaxIdleConnsPerHost: 10,               // Max number of idle connections per host
			IdleConnTimeout:     90 * time.Second, // Timeout for idle connections
			DisableKeepAlives:   false,            // Keep connections alive
		},
		Timeout: c.Config.Timeout,
	}

	c.httpClient = hc
	return c
}
