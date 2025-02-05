package client

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

// parseResponseBody returns the string representation of the response. From
// here, other formatting can be applied.
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

// setResponseHeaders sets the headers for the response on the API client.
func (api *ApiClient) setResponseHeaders(response *http.Response) {
	api.addHeaders(response.Header)
}

func buildQueryParams(query []ApiQuery) string {
	queryBuffer := bytes.NewBufferString("?")

	for _, q := range query {
		queryBuffer.WriteString(q.Key + "=" + url.QueryEscape(q.Value) + "&")
	}

	queryString := queryBuffer.String()

	return queryString[:len(queryString)-1] // remove the trailing '&'
}

// addHeaders adds headers to the provided http.Header.
// When the apiHeaders parameter is nil, it means that the headers are to be
// set on the ApiClient's ResponseHeaders field. This effectively serves as
// a way to set header values for requests and responses.
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

// resetDebugInfo resets the debug info built for a previous request-response
// cycle. This ensures that the information is always up-to-date.
func (api *ApiClient) resetDebugInfo() {
	api.debugInfo.Reset()
}

// setDebugInfo sets the debug info for a request-response cycle. It saves the
// information of both client and server during the conversation.
func (api *ApiClient) setDebugInfo(
	request *http.Request, response *http.Response,
) error {
	api.debugInfo.WriteString("API Debug Info\n===============\n\n")

	reqOut, reqOutErr := httputil.DumpRequest(request, true)
	if reqOutErr != nil {
		return reqOutErr
	}

	api.debugInfo.WriteString("Client Side\n============\n")
	api.debugInfo.WriteString(string(reqOut))

	resOut, resOutErr := httputil.DumpResponse(response, true)
	if resOutErr != nil {
		return resOutErr
	}

	api.debugInfo.WriteString("Server Side\n============\n")
	api.debugInfo.WriteString(string(resOut))

	return nil
}

// GetDebugInfo returns the debugged data collected during a request-response
// cycle.
func (api *ApiClient) GetDebugInfo() string {
	return api.debugInfo.String()
}

// requestHandler handles the request.
func (api *ApiClient) requestHandler(data *requestData) (
	*http.Response, error,
) {
	api.resetDebugInfo() // reset the debug info

	client := &http.Client{Timeout: api.Timeout}
	var queryString string

	if data.query != nil {
		queryString = buildQueryParams(data.query)
	}

	request, err := http.NewRequest(data.method, data.url+queryString, data.body)
	if err != nil {
		return nil, err
	}

	api.addHeaders(request.Header, data.headers...)
	api.addHeaders(
		request.Header, ApiHeader{Key: "User-Agent", Value: "httpClient v0.1"},
	)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if api.Debug {
		err = api.setDebugInfo(request, response)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

// responseHandler handles the response received from the server. It sets the
// status code and relevant response headers for client.
func (api *ApiClient) responseHandler(response *http.Response) {
	responseBody, err := parseResponseBody(response)
	if err != nil {
		api.Error = err
		return
	}

	api.StatusCode = response.StatusCode
	api.Body = responseBody
	api.setResponseHeaders(response)
}

// actionHandler handles the HTTP action to be performed.
func (api *ApiClient) actionHandler(data *requestData) {
	response, err := api.requestHandler(data)
	if err != nil {
		api.Error = err
		return
	}

	api.responseHandler(response)
}

// ResponseToMap takes the JSON response body and returns a map type for easy
// access and retrievals. This will fail if the Body is JSON unencodable.
func (api *ApiClient) ResponseToMap(m interface{}) (err error) {
	return responseToOther(m, api.Body)
}

// ResponseToStruct takes the JSON response body and returns a struct type for
// easy access and retrievals. This will fail if the Body is JSON unencodable.
func (api *ApiClient) ResponseToStruct(v interface{}) (err error) {
	return responseToOther(v, api.Body)
}

// responseToOther converts the API response body to the requested interface.
func responseToOther(output interface{}, responseBody string) (err error) {
	err = json.Unmarshal([]byte(responseBody), &output)
	if err != nil {
		return err
	}

	return nil
}
