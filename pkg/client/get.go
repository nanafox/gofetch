package client

import (
	"net/http"
)

// Get performs an API GET request.
func (api *ApiClient) Get(url string, query []ApiQuery, headers ...ApiHeader) {
	response, err := handleRequest(url, query, headers...)
	if err != nil {
		api.Error = err
		return
	}

	responseBody, err := parseResponseBody(response)
	if err != nil {
		api.Error = err
		return
	}

	api.StatusCode = response.StatusCode
	api.Body = responseBody
	api.setHeaders(response)
}

// handleRequest handles the GET request.
func handleRequest(
	url string, query []ApiQuery, headers ...ApiHeader,
) (*http.Response, error) {
	client := &http.Client{}

	queryString := buildQueryParams(query)

	request, err := http.NewRequest("GET", url+queryString, nil)
	if err != nil {
		return nil, err
	}

	addHeaders(request.Header, headers)

	return client.Do(request)
}
