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
func parseResponseBody(res *http.Response) (string, error) {
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// setResponseHeaders sets the headers for the response on the client.
func (c *Client) setResponseHeaders(res *http.Response) {
	c.addHeaders(res.Header)
}

func buildQueryParams(q []Query) string {
	buf := bytes.NewBufferString("?")

	for _, item := range q {
		buf.WriteString(item.Key + "=" + url.QueryEscape(item.Value) + "&")
	}

	s := buf.String()

	return s[:len(s)-1] // remove the trailing '&'
}

// addHeaders adds headers to the provided http.Header.
func (c *Client) addHeaders(h http.Header, hdrs ...Header) {
	if hdrs == nil {
		c.ResponseHeaders = make(map[string]string)

		for k, v := range h {
			c.ResponseHeaders[k] = strings.Join(v, " ")
		}
	} else {
		for _, hdr := range hdrs {
			h.Add(hdr.Key, hdr.Value)
		}
	}
}

// resetDebugInfo resets the debug info built for a previous request-response cycle.
func (c *Client) resetDebugInfo() {
	c.debugInfo.Reset()
}

// setDebugInfo sets the debug info for a request-response cycle.
func (c *Client) setDebugInfo(req *http.Request, res *http.Response) error {
	c.debugInfo.WriteString("API Debug Info\n===============\n\n")

	reqOut, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}

	c.debugInfo.WriteString("Client Side\n============\n")
	c.debugInfo.WriteString(string(reqOut))

	resOut, err := httputil.DumpResponse(res, true)
	if err != nil {
		return err
	}

	c.debugInfo.WriteString("Server Side\n============\n")
	c.debugInfo.WriteString(string(resOut))

	return nil
}

// GetDebugInfo returns the debugged data collected during a request-response cycle.
func (c *Client) GetDebugInfo() string {
	return c.debugInfo.String()
}

// requestHandler handles the request.
func (c *Client) requestHandler(data *requestData) (*http.Response, error) {
	c.resetDebugInfo() // reset the debug info

	// Use the optimized HTTP client that was created when the Client was instantiated
	var qs string
	if data.query != nil {
		qs = buildQueryParams(data.query)
	}

	// Build the request
	req, err := http.NewRequest(data.method, data.url+qs, data.body)
	if err != nil {
		return nil, err
	}

	// Set headers
	c.addHeaders(req.Header, data.headers...)
	c.addHeaders(req.Header, Header{Key: "User-Agent", Value: "httpClient v0.1"})

	// ensure the timeout is always the same as the configured one in case the
	// user changes it.
	c.httpClient.Timeout = c.Config.Timeout

	// Execute the request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Config.Debug {
		err = c.setDebugInfo(req, res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

// responseHandler handles the response received from the server.
func (c *Client) responseHandler(res *http.Response) {
	body, err := parseResponseBody(res)
	if err != nil {
		c.Error = err
		return
	}

	c.StatusCode = res.StatusCode
	c.Body = body
	c.setResponseHeaders(res)
}

// actionHandler handles the HTTP action to be performed.
func (c *Client) actionHandler(data *requestData) {
	res, err := c.requestHandler(data)
	if err != nil {
		c.Error = err
		return
	}

	c.responseHandler(res)
}

// ResponseToMap takes the JSON response body and returns a map type for easy access.
func (c *Client) ResponseToMap(m interface{}) (err error) {
	return responseToOther(m, c.Body)
}

// ResponseToStruct takes the JSON response body and returns a struct type for easy access.
func (c *Client) ResponseToStruct(v interface{}) (err error) {
	return responseToOther(v, c.Body)
}

// responseToOther converts the client response body to the requested interface.
func responseToOther(out interface{}, body string) (err error) {
	err = json.Unmarshal([]byte(body), &out)
	if err != nil {
		return err
	}

	return nil
}
