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
func (api *Client) Get(url string, query []Query, headers ...Header) {
	api.Do("GET", url, query, nil, headers...)
}

// Post performs an API POST request.
func (api *Client) Post(
	url string, query []Query, body io.Reader, headers ...Header,
) {
	api.Do("POST", url, query, body, headers...)
}

// Put performs an API PUT request.
func (api *Client) Put(
	url string, query []Query, body io.Reader, headers ...Header,
) {
	api.Do("PUT", url, query, body, headers...)
}

// Delete performs an API DELETE request.
func (api *Client) Delete(
	url string, query []Query, body io.Reader, headers ...Header,
) {
	api.Do("DELETE", url, query, body, headers...)
}
