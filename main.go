package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nanafox/simple-http-client/pkg/client"
)

func main() {
	apiClient := client.ApiClient{Debug: true}

	headers := []client.ApiHeader{
		{Key: "Accept", Value: "application/json"},
	}

	url := os.Getenv("API_URL")
	if url == "" {
		url = "https://httpbin.org/get"
	}

	queryParams := []client.ApiQuery{
		{Key: "name", Value: "awesome"},
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
}
