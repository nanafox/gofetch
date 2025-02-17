# gofetch - A Lightweight and Efficient HTTP Client for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/nanafox/gofetch.svg)](https://pkg.go.dev/github.com/nanafox/gofetch)

`gofetch` is a lightweight, high-performance HTTP client designed to simplify
making HTTP requests in Go. Built on top of Go’s standard `net/http` package,
`gofetch` provides a clean, intuitive API with minimal configuration while
offering powerful features like request timeouts, debugging, and structured
error handling.

This project started as a personal exploration of HTTP clients in Go. Over time,
it evolved into a practical tool, which I continue to maintain and improve.

---

## 🚀 Features

- **Simplified API** – Easy-to-use methods for common HTTP operations.
- **Minimal Dependencies** – Uses only Go's standard library.
- **Built-in Debugging** – Toggle request/response debugging on or off.
- **Customizable Configuration** – Set timeouts, headers, and other options.
- **Convenient Response Handling** – Access response status, body, and errors easily.
- **Flexible Query Parameters & Headers** – Use structured input for cleaner code.

---

## 📦 Installation

```sh
 go get -u github.com/nanafox/gofetch
```

Then import it into your Go project:

```go
import "github.com/nanafox/gofetch"
```

---

## 🔧 Configuration

You can create a `gofetch` client with optional configurations:

```go
client := gofetch.New(gofetch.Config{
    Timeout: 5 * time.Second,
    Debug: true,
})
```

**NOTE:**: If no configuration is provided, `gofetch` uses sensible defaults:

- **Timeout**: `500ms`
- **Debug**: `false`

You can also update the configuration later:

```go
client.Config.Debug = true  // Enable debugging
client.Config.Timeout = 10 * time.Second  // Change timeout
```

---

## 🌍 Making HTTP Requests

gofetch provides convenient methods for making HTTP requests:

### 1️⃣ GET Request

#### Get Method Signature

```go
package gofetch

// Get handles HTTP Get requests
func (api *Client) Get(url string, query []Query, headers ...Header)
```

```go
package main

import (
    "fmt"
    "log"
    "github.com/nanafox/gofetch"
)

func main() {
    client := gofetch.New()
    url := "https://httpbin.org/get"
 
    client.Get(url, nil)
 
    if client.Error != nil {
        log.Fatal(client.Error)
    }
 
    fmt.Println("Response:", client.Body)
}
```

#### ✅ With Query Parameters

```go
queryParams := []gofetch.Query{
    {Key: "search", Value: "gofetch"},
    {Key: "page", Value: "1"},
}
client.Get("https://httpbin.org/get", queryParams)
```

#### ✅ With Custom Headers

```go
headers := []gofetch.Header{
    {Key: "Authorization", Value: "Bearer token123"},
    {Key: "User-Agent", Value: "gofetch-client"},
}
client.Get("https://httpbin.org/get", nil, headers...)
```

---

### 2️⃣ POST Request

#### Post Method Signature

```go
package gofetch

// Get handles HTTP POST requests
func (api *Client) Post(url string, query []Query, body io.Reader, headers ...Header)
```

```go
body := `{ "username": "johndoe", "password": "secret" }`
headers := []gofetch.Header{{Key: "Content-Type", Value: "application/json"}}

client.Post("https://httpbin.org/post", nil, body, headers...)
```

---

### 3️⃣ PUT Request

#### Put Method Signature

```go
package gofetch

// Get handles HTTP PUT requests
func (api *Client) Put(url string, query []Query, body io.Reader, headers ...Header)
```

```go
body := `{ "name": "Updated Name" }`
client.Put("https://httpbin.org/put", nil, body)
```

---

### 4️⃣ DELETE Request

#### Delete Method Signature

```go
package gofetch

// Delete performs an API DELETE request.
func (api *Client) Delete(url string, query []Query, body io.Reader, headers ...Header)
```

#### Example

```go
client.Delete("https://httpbin.org/delete", nil, nil)
```

---

## 🛠 Using `Do()` for Flexible Requests

If you need more control, `Do()` lets you specify the HTTP method:

```go
client.Do("PATCH", "https://httpbin.org/patch", nil, `{"updated": true}`)
```

---

## 🕵️ Debugging & Error Handling

### ✅ Enable Debug Mode

Set `Debug: true` in the configuration to print request/response details:

```go
client := gofetch.New(gofetch.Config{Debug: true})
client.Get("https://httpbin.org/get", nil)
fmt.Println(client.GetDebugInfo())
```

### ✅ Handling Errors

If an error occurs, `client.Error` will contain details:

```go
if client.Error != nil {
    log.Fatalf("Request failed: %v", client.Error)
}
```

---

## 🔥 Why Use gofetch?

- 🏎 **Faster development** – Write clean, readable HTTP requests quickly.
- 🎯 **Minimalist & Lightweight** – No external dependencies.
- 🔍 **Debugging Friendly** – Get request/response details easily.
- 💡 **Customizable** – Supports timeouts, headers, query params, and more.

---

## 📜 License

`gofetch` is open-source and licensed under the MIT License.

---

## ❤️ Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

---

## 📫 Contact

- Author: **Maxwell Nana Forson (theLazyProgrammer)**
- X: [@_nanafox](https://x.com/_nanafox)
- GitHub: [@nanafox](https://github.com/nanafox)
- Website: [www.mnforson.live](https://www.mnforson.live)

---

Check the Reference

[![Go Reference](https://pkg.go.dev/badge/github.com/nanafox/gofetch.svg)](https://pkg.go.dev/github.com/nanafox/gofetch)

Happy coding! 🚀
