package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type sampleStruct struct {
	Login string `json:"age"`
	Email string `json:"email"`
	ID    int    `json:"id"`
}

func TestWriteJSON(t *testing.T) {
	// Create a mock HTTP response recorder
	w := httptest.NewRecorder()

	u := sampleStruct{
		Login: "tony",
		Email: "tony@tester",
		ID:    1,
	}

	// Call the WriteJSON function with the test data
	err := WriteJSON(w, http.StatusOK, u)
	if err != nil {
		t.Errorf("WriteJSON returned an error: %v", err)
	}

	// Verify the HTTP status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	// Verify the Content-Type header
	contentType := w.Header().Get("Content-Type")
	expectedContentType := "application/json"
	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type header to be %s, but got %s", expectedContentType, contentType)
	}

	// Parse the response body and compare it to the test data
	var responseJSON sampleStruct
	err = json.Unmarshal(w.Body.Bytes(), &responseJSON)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Compare the response data to the test data
	if !reflect.DeepEqual(responseJSON, u) {
		t.Errorf("Response data does not match the test data. Expected %+v, but got %+v", u, responseJSON)
	}
}

func TestReadJSON(t *testing.T) {
	// Create a mock HTTP request with a JSON payload
	requestBody := `{"age": "tony", "email": "tony@tester", "id": 1}`
	req, err := http.NewRequest("POST", "/example", strings.NewReader(requestBody))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Create a mock HTTP response recorder
	w := httptest.NewRecorder()

	// Define a struct to hold the JSON data
	testData := sampleStruct{}

	// Call the ReadJSON function with the mock request and the data struct
	err = ReadJSON(w, req, &testData)

	// Check for errors
	if err != nil {
		t.Errorf("ReadJSON returned an error: %v", err)
	}

	// Verify that the data was correctly decoded from the request body
	expectedData := sampleStruct{
		Login: "tony",
		Email: "tony@tester",
		ID:    1,
	}

	if !reflect.DeepEqual(testData, expectedData) {
		t.Errorf("Decoded data does not match the expected data. Expected %+v, but got %+v", expectedData, testData)
	}

	// Verify that the request body was read completely (should be EOF)
	buf := make([]byte, 1024)
	_, readErr := req.Body.Read(buf)
	if !errors.Is(readErr, io.EOF) {
		t.Errorf("Expected request body to be fully consumed (EOF), but got error: %v", readErr)
	}
}
