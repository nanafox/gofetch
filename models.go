package gofetch

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

// Client is a simple API client that makes HTTP requests.
// It contains the status code, response body, error information,
// response headers, debug information, configuration, and an HTTP client.
type Client struct {
	StatusCode      int
	Body            string
	Error           error
	ResponseHeaders map[string]string
	debugInfo       bytes.Buffer
	Config          Config
	httpClient      *http.Client
}

// Header defines the structure to model the API header key and values.
type Header struct {
	Key   string
	Value string
}

// Query defines a way to easily specify query parameters
type Query struct {
	Key   string
	Value string
}

// requestData is a struct that holds the data for an API request.
type requestData struct {
	method  string
	url     string
	query   []Query
	body    io.Reader
	headers []Header
}

// Config is used to store the configuration for the client
type Config struct {
	Timeout time.Duration
	Debug   bool
}
