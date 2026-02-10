package gofetch

import (
	"io"
)

// Do performs an API request with the specified HTTP method.
func (client *Client) Do(
	method, url string, query []Query, body io.Reader, headers ...Header,
) {
	data := &requestData{
		method: method, url: url, query: query, headers: headers, body: body,
	}

	client.actionHandler(data)
}

// Get performs an API GET request.
func (client *Client) Get(url string, query []Query, headers ...Header) {
	client.Do("GET", url, query, nil, headers...)
}

// Post performs an API POST request.
func (client *Client) Post(
	url string, query []Query, body io.Reader, headers ...Header,
) {
	client.Do("POST", url, query, body, headers...)
}

// Put performs an API PUT request.
func (client *Client) Put(
	url string, query []Query, body io.Reader, headers ...Header,
) {
	client.Do("PUT", url, query, body, headers...)
}

// Delete performs an API DELETE request.
func (client *Client) Delete(
	url string, query []Query, body io.Reader, headers ...Header,
) {
	client.Do("DELETE", url, query, body, headers...)
}
