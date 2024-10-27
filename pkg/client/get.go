package client

import (
	"net/http"
)

// Get performs an API GET request.
func (api *ApiClient) Get(url string, query []ApiQuery, headers ...ApiHeader) {
	response, err := api.handleRequest(url, query, headers...)
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
	api.setResponseHeaders(response)
}

// handleRequest handles the GET request.
func (api *ApiClient) handleRequest(
	url string, query []ApiQuery, headers ...ApiHeader,
) (*http.Response, error) {
	api.resetDebugInfo() // reset the debug info

	client := &http.Client{}

	queryString := buildQueryParams(query)

	request, err := http.NewRequest("GET", url+queryString, nil)
	if err != nil {
		return nil, err
	}

	api.addHeaders(request.Header, headers...)
	api.addHeaders(
		request.Header, ApiHeader{Key: "User-Agent", Value: "httpClient v0.1"},
	)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if api.Debug == true {
		err = api.setDebugInfo(request, response)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}
