package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

// parseResponseBody returns the string representation of the response. From
// here, other formatting can be applied.
func parseResponseBody(response *http.Response) (string, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("An error occurred while reading the response: %v\n", err)
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

// setResponseHeaders sets the headers for the response on the API client.
func (api *ApiClient) setResponseHeaders(response *http.Response) {
	api.addHeaders(response.Header)
}

func buildQueryParams(query []ApiQuery) string {
	queryBuffer := bytes.NewBufferString("?")

	for _, q := range query {
		queryBuffer.WriteString(fmt.Sprintf("%v=%v&", q.Key, q.Value))
	}

	queryString := queryBuffer.String()

	return queryString[:len(queryString)-1] // remove the trailing '&'
}

// addHeaders adds headers to the provided http.Header.
// When the apiHeaders parameter is nil, it means that the headers are to be
// set on the ApiClient's ResponseHeaders field. This effectively serves as a way to
// set header values for requests and responses.
func (api *ApiClient) addHeaders(header http.Header, apiHeaders ...ApiHeader) {
	if apiHeaders == nil {
		api.ResponseHeaders = make(map[string]string)

		for key, value := range header {
			api.ResponseHeaders[key] = strings.Join(value, " ")
		}
	} else {
		for _, h := range apiHeaders {
			header.Add(h.Key, h.Value)
		}
	}
}
