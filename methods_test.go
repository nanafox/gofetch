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
var srv *httptest.Server

// Global debugging flag
var debugEnabled bool

// Initialize the global test server and set up debugging
func TestMain(m *testing.M) {
	// Initialize the test server with different response behaviors
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK) // Simulate a 200 OK response for GET
			data = []byte(`{"message": "GET request successful"}`)
		case "POST":
			w.WriteHeader(http.StatusCreated) // Simulate a 201 Created response for POST
			data = []byte(`{"message": "POST request successful"}`)
		case "PUT":
			w.WriteHeader(http.StatusOK)
			data = []byte(`{"message": "PUT request successful"}`)
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
			data = []byte(nil)
		}

		// Write the response data
		_, err := w.Write(data)
		if err != nil {
			return
		}
	}))

	// Run the tests
	code := m.Run()

	// Cleanup: Stop the server after all tests
	srv.Close()

	// Exit with the result code from running tests
	if code != 0 {
		panic("Tests failed")
	}
}

// Example Test Cases Using the Global Test Server
func TestClientDoWithMockServer(t *testing.T) {
	// Flag to enable debugging for all tests
	debugEnabled = true

	tests := []struct {
		method       string
		url          string
		expectedCode int
		expectedData map[string]string
	}{
		{
			"GET",
			srv.URL + "/test",
			http.StatusOK,
			map[string]string{"message": "GET request successful"},
		},
		{
			"POST",
			srv.URL + "/test",
			http.StatusCreated,
			map[string]string{"message": "POST request successful"},
		},
		{
			"PUT",
			srv.URL + "/test",
			http.StatusOK,
			map[string]string{"message": "PUT request successful"},
		},
		{
			"DELETE",
			srv.URL + "/test",
			http.StatusNoContent,
			nil,
		},
	}

	// Run each test
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			// Create a new client and set debugging if enabled
			c := New(Config{})
			c.Config.Debug = debugEnabled

			var body *bytes.Buffer
			if tt.method == "POST" || tt.method == "PUT" {
				body = bytes.NewBufferString("test data")
			}

			// Call the client method based on the HTTP method
			switch tt.method {
			case "GET":
				c.Do("GET", tt.url, nil, nil)
			case "POST":
				c.Do("POST", tt.url, nil, body)
			case "PUT":
				c.Do("PUT", tt.url, nil, body)
			case "DELETE":
				c.Do("DELETE", tt.url, nil, nil)
			}

			// Assert the status code matches the expected result
			assert.Equal(t, tt.expectedCode, c.StatusCode, "Status code mismatch")

			// Unmarshal the response data and assert the expected response
			var res map[string]string
			err := c.ResponseToMap(&res)
			if err != nil {
				return
			}
			if c.Error != nil {
				t.Fatalf("Error while using ResponseToMap: %v", c.Error)
			}

			// Assert that the message in the response matches the expected message
			assert.Equal(
				t,
				tt.expectedData["message"],
				res["message"],
				"Response data mismatch",
			)

			if c.GetDebugInfo() == "" && c.Config.Debug == true {
				t.Fatal("GetDebugInfo() must return an empty string when debug is set to true")
			}
		})
	}
}

// TestClientGet tests the Get method and validates query parameters
func TestClientGet(t *testing.T) {
	// Enable debugging for this test
	debugEnabled = true

	// Start the test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the query parameters as JSON
		q := r.URL.Query()
		res := map[string]interface{}{
			"message": "GET request successful",
			"query":   q,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Respond with query parameters received
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			return
		}
	}))
	defer srv.Close()

	c := New(Config{Debug: debugEnabled})

	// Example query parameters
	query := []Query{
		{"param1", "value1"},
		{"param2", "value2"},
	}

	// Make a GET request with query parameters
	c.Get(srv.URL, query)

	// Assert the status code is what we expect for a GET request
	assert.Equal(t, http.StatusOK, c.StatusCode, "Expected status code 200 for GET request")

	// Use ResponseToMap to decode the response
	var res map[string]interface{}
	err := c.ResponseToMap(&res)
	if err != nil {
		return
	}
	if c.Error != nil {
		t.Fatalf("Error while using ResponseToMap: %v", c.Error)
	}

	// Assert the response message is correct
	assert.Equal(
		t,
		"GET request successful",
		res["message"],
		"Expected response message for GET request",
	)

	// Assert that the query parameters were received correctly
	qMap := res["query"].(map[string]interface{})

	assert.Equal(
		t,
		"value1",
		qMap["param1"].([]interface{})[0],
		"Expected param1 to be 'value1'",
	)
	assert.Equal(
		t,
		"value2",
		qMap["param2"].([]interface{})[0],
		"Expected param2 to be 'value2'",
	)
}

// TestClientPost tests the Post method of the Client
func TestClientPost(t *testing.T) {
	// Enable debugging for this test
	debugEnabled = true

	// Start the test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Respond with the received request body
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return
		}
	}))
	defer srv.Close()

	c := New(Config{Debug: debugEnabled})

	// Prepare POST request body
	body := bytes.NewBufferString(`{"key":"value"}`)

	// Make a POST request
	c.Post(srv.URL, nil, body)

	// Assert the status code is what we expect for a POST request
	assert.Equal(
		t,
		http.StatusCreated,
		c.StatusCode,
		"Expected status code 201 for POST request",
	)

	// Use ResponseToMap to decode the response
	var res map[string]string
	err := c.ResponseToMap(&res)
	if err != nil {
		return
	}
	if c.Error != nil {
		t.Fatalf("Error while using ResponseToMap: %v", c.Error)
	}

	// Assert the response message is correct
	assert.Equal(t, "value", res["key"], "Expected key to have value 'value'")
}

// TestClientPut tests the Put method of the Client
func TestClientPut(t *testing.T) {
	// Enable debugging for this test
	debugEnabled = true

	// Start the test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Respond with the received request body (optional, as we might not send back content for 204)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return
		}
	}))
	defer srv.Close()

	c := New(Config{Debug: debugEnabled})

	// Prepare the PUT request body
	body := bytes.NewBufferString(`{"key":"updated value"}`)

	// Make a PUT request
	c.Put(srv.URL, nil, body)

	// Assert the status code is what we expect for a PUT request
	assert.Equal(
		t,
		http.StatusOK,
		c.StatusCode,
		"Expected status code 204 for PUT request",
	)

	// Use ResponseToMap to decode the response (optional, as the response might be empty for 204)
	var res map[string]string
	err := c.ResponseToMap(&res)
	if err != nil {
		t.Fatalf("Error while using ResponseToMap: %v", err)
	}
	if c.Error != nil {
		t.Fatalf("Error while using ResponseToMap: %v", c.Error)
	}

	// Assert the response message is correct (optional)
	assert.Equal(t, "updated value", res["key"], "Expected key to be 'updated value'")
}

// TestClientDelete tests the Delete method of the Client
func TestClientDelete(t *testing.T) {
	// Enable debugging for this test
	debugEnabled = true

	// Start the test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the method used for the delete operation
		data := map[string]string{
			"message": "DELETE request successful",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Respond with the message
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return
		}
	}))
	defer srv.Close()

	c := New(Config{Debug: debugEnabled})

	// Make a DELETE request
	c.Delete(srv.URL, nil, nil)

	// Assert the status code is what we expect for a GET request
	assert.Equal(t, http.StatusOK, c.StatusCode, "Expected status code 200 for DELETE request")

	// Use ResponseToMap to decode the response
	var res map[string]string
	err := c.ResponseToMap(&res)
	if err != nil {
		return
	}
	if c.Error != nil {
		t.Fatalf("Error while using ResponseToMap: %v", c.Error)
	}

	// Assert the response message is correct
	assert.Equal(
		t,
		"DELETE request successful",
		res["message"],
		"Expected response message for DELETE request",
	)
}
