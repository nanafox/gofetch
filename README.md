# A Simple Go HTTP Client

This package aims to simplify HTTP requests by wrapping around the
`net-http` package. This is a hobby project to explore HTTP clients in Go.

## Installation

1. Download the package
    ```bash
    go get -u github.com/nanafox/simple-http-client
    ```

2. Use it for your HTTP requests

## Features

- Supports setting request timeouts.
- Supports enabling and disabling debug information.

## Example Usages

### GET Requests

**Method Signature:**

```go
package client

func (api *ApiClient) Get(url string, query []ApiQuery, headers ...ApiHeader)
```

#### Example 1 - Simple API request, no query params or custom headers

For simple requests, the `query` and `headers` can be omitted if not needed.
Specify `nil` for `query` and nothing at all for the `headers`.

**Note: To enable debugging, set the `Debug` field to `true` on your instance**
For example: `apiClient.Debug = true`

```go
package main

import (
   "fmt"
   "github.com/nanafox/simple-http-client/pkg/client"
   "log"
)

func main() {
   apiClient := client.ApiClient{}
   url := "https://httpbin.org/get"

   apiClient.Get(url, nil)

   if apiClient.Error != nil {
      log.Fatal(apiClient.Error)
   }
   
   fmt.Printf("The status code was: %v\n", apiClient.StatusCode)
   fmt.Println(apiClient.Body)
}
```

#### Example 2 - Using Debugging Mode and a Timeout value

By default, debugging is turned off but there are times when you'd want to
peak a bit more under hood to see how the client (you) appeared to the server
and how the server responded. For this, set `Debug` field to `true`. After it
has been turned on, you can retrieve the debugged info with the `GetDebugInfo()`
method. Once you don't need debugging, simply set `Debug` to `false`.

Another default is that requests will hang indefinitely if the server is taking
too long to respond. There could be multiple factors to this which are mostly
out of your control. Setting the timeout for the API client ensures that you are
blocking the system indefinitely. Instead, you set a window of the expected time
you want respond to come in. If the response is not received by the time the
duration has passed, a timeout error is returned. Adjust this value as you want.

```go
package main

import (
   "fmt"
   "github.com/nanafox/simple-http-client/pkg/client"
   "log"
   "time"
)

func main() { 
   // Debug and Timeout can be set while creating an API Client instance
   apiClient := client.ApiClient{Debug: true, Timeout: time.Second * 5}

   // they can also be set individually after an instance has been created
   apiClient.Debug = true // enable debugging
   apiClient.Timeout = time.Second * 5 // Timeout after five seconds
   
   // Make a request
   url := "https://httpbin.org/get"
   apiClient.Get(url, nil)
	 
   // handle errors
   if apiClient.Error != nil {
      log.Fatal(apiClient.Error)
   }
	 
   // Print the debugging info
   fmt.Println(apiClient.GetDebugInfo())
}
```

#### Example 3 - Passing Query Parameters

There are times when the server allows you to filter the response it returns by
providing query parameters. The API client allows you to conveniently list all
the query parameters you want to use.

It provides the `ApiQuery` struct to ease things up a bit. Provide the parameters
you need as an array of `ApiQuery` objects. It's a simple key-value struct.
Specify the `Key` and `Value` as strings.

Check the code sample below for an example usage.

```go
// the usual setup omitted for brevity

func main() {
   apiClient := client.ApiClient{Timeout: time.Second * 5}

	 // Add query parameters: simply key-value pairs
   queryParams := []client.ApiQuery{
    {Key: "q", Value: "HTTP Clients"},
    {Key: "page", Value: "3"},
    {Key: "page_size", Value: "10"},
    // add as many more as you need
  }

   url := "https://httpbin.org/get"
	 
   // The API will build the full URL with the query parameters, for example
   // the resultant URL of this request will be
   // GET https://httpbin.ord/get?q=HTTP+Clients&page=3&page_size=10
   
  // make the request
   apiClient.Get(url, queryParams)

  // response omitted for brevity
}
```

#### Example 4 - Setting Request Headers

Setting headers is a common thing for HTTP requests. To set headers, you can
provide them as an array of `ApiHeader` objects. The `ApiHeader` struct is a
simple key-value struct. Provide the `Key` and `Value` as strings.

```go
// the usual setup omitted for brevity

func main() {
   apiClient := client.ApiClient{Timeout: time.Second * 5}

   // Add headers: simply key-value pairs
   headers := []client.ApiHeader{
    {Key: "Authorization" Value: "Bearer <your-token-here>"},
    {Key: "Content-Type", Value: "application/json"},
    // add as many more as you need
  }

    url := "https://httpbin.org/get"
    
    // make the request
    apiClient.Get(url, nil, headers...)

    // response omitted for brevity
}
```
