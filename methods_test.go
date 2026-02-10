package gofetch

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Global test server
var testServer *httptest.Server

// Global debugging flag
var DebuggingEnabled bool

// Initialize the global test server and set up debugging
func TestMain(m *testing.M) {
	// Initialize the test server with different response behaviors
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var responseData []byte
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK) // Simulate a 200 OK response for GET
			responseData = []byte(`{"message": "GET request successful"}`)
		case "POST":
			w.WriteHeader(http.StatusCreated) // Simulate a 201 Created response for POST
			responseData = []byte(`{"message": "POST request successful"}`)
		case "PUT":
			w.WriteHeader(http.StatusOK)
			responseData = []byte(`{"message": "PUT request successful"}`)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
			responseData = []byte(nil)
		}

		// Write the response data
		_, err := w.Write(responseData)
		if err != nil {
			return
		}
	}))

	// Run the tests
	code := m.Run()

	// Cleanup: Stop the server after all tests
	testServer.Close()

	// Exit with the result code from running tests
	if code != 0 {
		panic("Tests failed")
	}
}

// Example Test Cases Using the Global Test Server
func TestClient_DoWithMockServer(t *testing.T) {
	// Flag to enable debugging for all tests
	DebuggingEnabled = true

	tests := []struct {
		method       string
		url          string
		expectedCode int
		expectedData map[string]string
	}{
		{
			"GET",
			testServer.URL + "/test",
			http.StatusOK,
			map[string]string{"message": "GET request successful"},
		},
		{
			"POST",
			testServer.URL + "/test",
			http.StatusCreated,
			map[string]string{"message": "POST request successful"},
		},
		{
			"PUT",
			testServer.URL + "/test",
			http.StatusOK,
			map[string]string{"message": "PUT request successful"},
		},
		{
			"DELETE",
			testServer.URL + "/test",
			http.StatusNoContent,
			nil,
		},
	}

	// Run each test
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			// Create a new client and set debugging if enabled
			client := New()
			client.Config.Debug = DebuggingEnabled

			var reqBody *bytes.Buffer
			if tt.method == "POST" || tt.method == "PUT" {
				reqBody = bytes.NewBufferString("test data")
			}

			// Call the client method based on the HTTP method
			switch tt.method {
			case "GET":
				client.Do("GET", tt.url, nil, nil)
			case "POST":
				client.Do("POST", tt.url, nil, reqBody)
			case "PUT":
				client.Do("PUT", tt.url, nil, reqBody)
			case "DELETE":
				client.Do("DELETE", tt.url, nil, nil)
			}

			// Assert the status code matches the expected result
			assert.Equal(t, tt.expectedCode, client.StatusCode, "Status code mismatch")

			// Unmarshal the response data and assert the expected response
			var response map[string]string
			err := client.ResponseToMap(&response)
			if err != nil {
				return
			}
			if client.Error != nil {
				t.Fatalf("Error while using ResponseToMap: %v", client.Error)
			}

			// Assert that the message in the response matches the expected message
			assert.Equal(
				t,
				tt.expectedData["message"],
				response["message"],
				"Response data mismatch",
			)

			if client.GetDebugInfo() == "" && client.Config.Debug == true {
				t.Fatal("GetDebugInfo() must return an empty string when debug is set to true")
			}
		})
	}
}

// TestClient_Get tests the Get method and validates query parameters
func TestClient_Get(t *testing.T) {
	// Enable debugging for this test
	DebuggingEnabled = true

	// Start the test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the query parameters as JSON
		queryParams := r.URL.Query()
		responseData := map[string]interface{}{
			"message": "GET request successful",
			"query":   queryParams,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Respond with query parameters received
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	client := New()
	client.Config.Debug = DebuggingEnabled

	// Example query parameters
	query := []Query{
		{"param1", "value1"},
		{"param2", "value2"},
	}

	// Make a GET request with query parameters
	client.Get(server.URL, query)

	// Assert the status code is what we expect for a GET request
	assert.Equal(t, http.StatusOK, client.StatusCode, "Expected status code 200 for GET request")

	// Use ResponseToMap to decode the response
	var response map[string]interface{}
	err := client.ResponseToMap(&response)
	if err != nil {
		return
	}
	if client.Error != nil {
		t.Fatalf("Error while using ResponseToMap: %v", client.Error)
	}

	// Assert the response message is correct
	assert.Equal(
		t,
		"GET request successful",
		response["message"],
		"Expected response message for GET request",
	)

	// Assert that the query parameters were received correctly
	queryMap := response["query"].(map[string]interface{})

	assert.Equal(
		t,
		"value1",
		queryMap["param1"].([]interface{})[0],
		"Expected param1 to be 'value1'",
	)
	assert.Equal(
		t,
		"value2",
		queryMap["param2"].([]interface{})[0],
		"Expected param2 to be 'value2'",
	)
}

// TestClient_Post tests the Post method of the Client
func TestClient_Post(t *testing.T) {
	// Enable debugging for this test
	DebuggingEnabled = true

	// Start the test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var responseData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&responseData); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Respond with the received request body
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	client := New()
	client.Config.Debug = DebuggingEnabled

	// Prepare POST request body
	reqBody := bytes.NewBufferString(`{"key":"value"}`)

	// Make a POST request
	client.Post(server.URL, nil, reqBody)

	// Assert the status code is what we expect for a POST request
	assert.Equal(
		t,
		http.StatusCreated,
		client.StatusCode,
		"Expected status code 201 for POST request",
	)

	// Use ResponseToMap to decode the response
	var response map[string]string
	err := client.ResponseToMap(&response)
	if err != nil {
		return
	}
	if client.Error != nil {
		t.Fatalf("Error while using ResponseToMap: %v", client.Error)
	}

	// Assert the response message is correct
	assert.Equal(t, "value", response["key"], "Expected key to have value 'value'")
}

// TestClient_Put tests the Put method of the Client
func TestClient_Put(t *testing.T) {
	// Enable debugging for this test
	DebuggingEnabled = true

	// Start the test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var responseData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&responseData); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Respond with the received request body (optional, as we might not send back content for 204)
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	client := New()
	client.Config.Debug = DebuggingEnabled

	// Prepare the PUT request body
	reqBody := bytes.NewBufferString(`{"key":"updated value"}`)

	// Make a PUT request
	client.Put(server.URL, nil, reqBody)

	// Assert the status code is what we expect for a PUT request
	assert.Equal(
		t,
		http.StatusOK,
		client.StatusCode,
		"Expected status code 204 for PUT request",
	)

	// Use ResponseToMap to decode the response (optional, as the response might be empty for 204)
	var response map[string]string
	err := client.ResponseToMap(&response)
	if err != nil {
		t.Fatalf("Error while using ResponseToMap: %v", err)
	}
	if client.Error != nil {
		t.Fatalf("Error while using ResponseToMap: %v", client.Error)
	}

	// Assert the response message is correct (optional)
	assert.Equal(t, "updated value", response["key"], "Expected key to be 'updated value'")
}

// TestClient_Delete tests the Delete method of the Client
func TestClient_Delete(t *testing.T) {
	// Enable debugging for this test
	DebuggingEnabled = true

	// Start the test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the method used for the delete operation
		responseData := map[string]string{
			"message": "DELETE request successful",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Respond with the message
		err := json.NewEncoder(w).Encode(responseData)
		if err != nil {
			return
		}
	}))
	defer server.Close()

	client := New()
	client.Config.Debug = DebuggingEnabled

	// Make a DELETE request
	client.Delete(server.URL, nil, nil)

	// Assert the status code is what we expect for a DELETE request
	assert.Equal(t, http.StatusOK, client.StatusCode, "Expected status code 200 for DELETE request")

	// Use ResponseToMap to decode the response
	var response map[string]string
	err := client.ResponseToMap(&response)
	if err != nil {
		return
	}
	if client.Error != nil {
		t.Fatalf("Error while using ResponseToMap: %v", client.Error)
	}

	// Assert the response message is correct
	assert.Equal(
		t,
		"DELETE request successful",
		response["message"],
		"Expected response message for DELETE request",
	)
}
