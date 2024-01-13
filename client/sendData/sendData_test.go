package senddata

import (
	"bytes"
	"devashishRaj/rpi_telemetry/server/handleError"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHttpPost(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, but got %s", r.Header.Get("Content-Type"))
		}

		// Verify request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		expectedJSON := []byte(`{"key":"value"}`) // Replace with your expected JSON
		if !bytes.Equal(body, expectedJSON) {
			t.Errorf("Expected request body to be %s, but got %s", expectedJSON, body)
		}

		// Send a dummy response
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		handleError.CheckError("error in writing dummy responser", err)
	}))

	// Close the server when the test is done
	defer server.Close()

	// Redirect log output to a buffer for assertions
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Call the HttpPost function with the mock server URL and JSON data
	URL := server.URL
	jsonData := []byte(`{"key":"value"}`) // Replace with your actual JSON data

	// Call the function
	HttpPost(URL, jsonData)

	// Verify received JSON data
	var receivedData map[string]interface{}
	err := json.Unmarshal(jsonData, &receivedData)
	if err != nil {
		t.Fatal(err)
	}

	// Replace with your expected data
	expectedData := map[string]interface{}{"key": "value"}
	if !reflect.DeepEqual(receivedData, expectedData) {
		t.Errorf("Expected received JSON data to be %v, but got %v", expectedData, receivedData)
	}
}
