package client

import (
	"bytes"
)

// ApiClient is a simple API client that makes HTTP requests fun.
type ApiClient struct {
	StatusCode      int
	Body            string
	Error           error
	ResponseHeaders map[string]string
	Debug           bool
	debugInfo       bytes.Buffer
}

// ApiHeader defines the structure to model API header key and values.
type ApiHeader struct {
	Key   string
	Value string
}

// ApiQuery defines a way to easily specify query parameters
type ApiQuery struct {
	Key   string
	Value any
}
