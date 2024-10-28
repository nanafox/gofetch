# A Simple Go HTTP Client

This package aims to simplify HTTP requests by wrapping around the
`net-http` package. This is a hobby project to see how the Go language.

## Installation

1. Download the package
    ```bash
    go get -u github.com/nanafox/simple-http-client
    ```

2. Use it for your HTTP requests

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
