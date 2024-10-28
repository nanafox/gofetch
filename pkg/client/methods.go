package client

import (
	"io"
)

// Do performs an API request with the specified HTTP method.
func (api *ApiClient) Do(
	method, url string, query []ApiQuery, body io.Reader, headers ...ApiHeader,
) {
	data := &requestData{
		method: method, url: url, query: query, headers: headers, body: body,
	}

	api.actionHandler(data)
}

// Get performs an API GET request.
func (api *ApiClient) Get(url string, query []ApiQuery, headers ...ApiHeader) {
	api.Do("GET", url, query, nil, headers...)
}

// Post performs an API POST request.
func (api *ApiClient) Post(
	url string, query []ApiQuery, body io.Reader, headers ...ApiHeader,
) {
	api.Do("POST", url, query, body, headers...)
}

// Put performs an API PUT request.
func (api *ApiClient) Put(
	url string, query []ApiQuery, body io.Reader, headers ...ApiHeader,
) {
	api.Do("PUT", url, query, body, headers...)
}

// Delete performs an API DELETE request.
func (api *ApiClient) Delete(
	url string, query []ApiQuery, body io.Reader, headers ...ApiHeader,
) {
	api.Do("DELETE", url, query, body, headers...)
}
