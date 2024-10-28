package client

import (
	"bytes"
	"io"
	"time"
)

// ApiClient is a simple API client that makes HTTP requests fun.
type ApiClient struct {
	StatusCode      int
	Body            string
	Error           error
	ResponseHeaders map[string]string
	Debug           bool
	debugInfo       bytes.Buffer
	Timeout         time.Duration
}

// ApiHeader defines the structure to model API header key and values.
type ApiHeader struct {
	Key   string
	Value string
}

// ApiQuery defines a way to easily specify query parameters
type ApiQuery struct {
	Key   string
	Value string
}

// requestData is a struct that holds the data for an API request.
type requestData struct {
	method  string
	url     string
	query   []ApiQuery
	body    io.Reader
	headers []ApiHeader
}
