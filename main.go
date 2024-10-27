package main

import (
	"fmt"

	"github.com/nanafox/simple-http-client/pkg/client"
)

func main() {
	apiClient := client.ApiClient{}

	headers := []client.ApiHeader{
		{Key: "Accept", Value: "application/json"},
	}

	url := "https://httpbin.org/get"

	queryParams := []client.ApiQuery{
		{Key: "name", Value: "awesome"},
		// Add as many as needed, it will handled automatically for you
	}

	apiClient.Get(url, queryParams, headers...)

	if apiClient.Error != nil {
		fmt.Println(apiClient.Error)
		return
	}

	fmt.Println(apiClient.StatusCode)
	fmt.Println(apiClient.Body)
}
