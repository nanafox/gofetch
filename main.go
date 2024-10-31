package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nanafox/simple-http-client/pkg/client"
)

// HttpBinResponse is a struct that represents the response from the
// httpbin.org API.
type HttpBinResponse struct {
	Args    map[string]string `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	Url     string            `json:"url"`
}

func main() {
	apiClient := client.ApiClient{Debug: true, Timeout: 5 * time.Second}

	headers := []client.ApiHeader{
		{Key: "Accept", Value: "application/json"},
	}

	url := os.Getenv("API_URL")
	if url == "" {
		url = "https://httpbin.org/get"
	}

	queryParams := []client.ApiQuery{
		{Key: "name", Value: "John Doe"},
		// Add as many as needed, it will be handled automatically for you
	}

	apiClient.Get(url, queryParams, headers...)

	if apiClient.Error != nil {
		log.Fatal(apiClient.Error)
	}

	// print the debug info for the request-response cycle
	fmt.Println(apiClient.GetDebugInfo())

	// the client keeps the response header information
	fmt.Println("\n\nResponse Headers")
	fmt.Println(apiClient.ResponseHeaders)

	// convert JSON response to Go map
	m, err := apiClient.ResponseToMap()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nResponseToMap: %+v\n", m)

	// convert JSON response to Go struct
	var httpBinResponse HttpBinResponse

	err = apiClient.ResponseToStruct(&httpBinResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nResponseToStruct: %+v\n", httpBinResponse)
}
