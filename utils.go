package gofetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// parseResponseBody returns the string representation of the response. From here, other formatting can be applied.
func parseResponseBody(response *http.Response) (string, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("an error occurred while reading the response: %v", err)
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

// setResponseHeaders sets the headers for the response on the client.
func (client *Client) setResponseHeaders(response *http.Response) {
	client.addHeaders(response.Header)
}

func buildQueryParams(query []Query) string {
	queryBuffer := bytes.NewBufferString("?")

	for _, q := range query {
		queryBuffer.WriteString(q.Key + "=" + url.QueryEscape(q.Value) + "&")
	}

	queryString := queryBuffer.String()

	return queryString[:len(queryString)-1] // remove the trailing '&'
}

// addHeaders adds headers to the provided http.Header.
func (client *Client) addHeaders(header http.Header, apiHeaders ...Header) {
	if apiHeaders == nil {
		client.ResponseHeaders = make(map[string]string)

		for key, value := range header {
			client.ResponseHeaders[key] = strings.Join(value, " ")
		}
	} else {
		for _, h := range apiHeaders {
			header.Add(h.Key, h.Value)
		}
	}
}

// resetDebugInfo resets the debug info built for a previous request-response cycle.
func (client *Client) resetDebugInfo() {
	client.debugInfo.Reset()
}

// setDebugInfo sets the debug info for a request-response cycle.
func (client *Client) setDebugInfo(request *http.Request, response *http.Response) error {
	client.debugInfo.WriteString("API Debug Info\n===============\n\n")

	reqOut, reqOutErr := httputil.DumpRequest(request, true)
	if reqOutErr != nil {
		return reqOutErr
	}

	client.debugInfo.WriteString("Client Side\n============\n")
	client.debugInfo.WriteString(string(reqOut))

	resOut, resOutErr := httputil.DumpResponse(response, true)
	if resOutErr != nil {
		return resOutErr
	}

	client.debugInfo.WriteString("Server Side\n============\n")
	client.debugInfo.WriteString(string(resOut))

	return nil
}

// GetDebugInfo returns the debugged data collected during a request-response cycle.
func (client *Client) GetDebugInfo() string {
	return client.debugInfo.String()
}

// requestHandler handles the request.
func (client *Client) requestHandler(data *requestData) (*http.Response, error) {
	client.resetDebugInfo() // reset the debug info

	// Use the optimized HTTP client that was created when the Client was instantiated
	var queryString string
	if data.query != nil {
		queryString = buildQueryParams(data.query)
	}

	// Build the request
	request, err := http.NewRequest(data.method, data.url+queryString, data.body)
	if err != nil {
		return nil, err
	}

	// Set headers
	client.addHeaders(request.Header, data.headers...)
	client.addHeaders(request.Header, Header{Key: "User-Agent", Value: "httpClient v0.1"})

	// ensure the timeout is always the same as the configured one in case the
	// user changes it.
	client.httpClient.Timeout = client.Config.Timeout

	// Execute the request
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if client.Config.Debug {
		err = client.setDebugInfo(request, response)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

// responseHandler handles the response received from the server.
func (client *Client) responseHandler(response *http.Response) {
	responseBody, err := parseResponseBody(response)
	if err != nil {
		client.Error = err
		return
	}

	client.StatusCode = response.StatusCode
	client.Body = responseBody
	client.setResponseHeaders(response)
}

// actionHandler handles the HTTP action to be performed.
func (client *Client) actionHandler(data *requestData) {
	response, err := client.requestHandler(data)
	if err != nil {
		client.Error = err
		return
	}

	client.responseHandler(response)
}

// ResponseToMap takes the JSON response body and returns a map type for easy access.
func (client *Client) ResponseToMap(m interface{}) (err error) {
	return responseToOther(m, client.Body)
}

// ResponseToStruct takes the JSON response body and returns a struct type for easy access.
func (client *Client) ResponseToStruct(v interface{}) (err error) {
	return responseToOther(v, client.Body)
}

// responseToOther converts the client response body to the requested interface.
func responseToOther(output interface{}, responseBody string) (err error) {
	err = json.Unmarshal([]byte(responseBody), &output)
	if err != nil {
		return err
	}

	return nil
}
