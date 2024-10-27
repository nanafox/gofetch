package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

// setHeaders sets the headers for the response on the API client.
func (api *ApiClient) setHeaders(response *http.Response) {
	addHeaders(response.Header, api.Headers)
}

func buildQueryParams(query []ApiQuery) string {
	queryBuffer := bytes.NewBufferString("?")

	for _, q := range query {
		queryBuffer.WriteString(fmt.Sprintf("%v=%v&", q.Key, q.Value))
	}

	return queryBuffer.String()
}

// addHeaders adds headers to the provided http.Header.
func addHeaders(header http.Header, apiHeaders []ApiHeader) {
	for _, h := range apiHeaders {
		header.Add(h.Key, h.Value)
	}
}
