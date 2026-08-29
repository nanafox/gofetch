package gofetch

import (
	"io"
)

// Do performs an API request with the specified HTTP method.
func (c *Client) Do(
	method, url string, query []Query, body io.Reader, headers ...Header,
) {
	data := &requestData{
		method: method, url: url, query: query, headers: headers, body: body,
	}

	c.actionHandler(data)
}

// Get performs an API GET request.
func (c *Client) Get(url string, query []Query, headers ...Header) {
	c.Do("GET", url, query, nil, headers...)
}

// Post performs an API POST request.
func (c *Client) Post(
	url string, query []Query, body io.Reader, headers ...Header,
) {
	c.Do("POST", url, query, body, headers...)
}

// Put performs an API PUT request.
func (c *Client) Put(
	url string, query []Query, body io.Reader, headers ...Header,
) {
	c.Do("PUT", url, query, body, headers...)
}

// Delete performs an API DELETE request.
func (c *Client) Delete(
	url string, query []Query, body io.Reader, headers ...Header,
) {
	c.Do("DELETE", url, query, body, headers...)
}
