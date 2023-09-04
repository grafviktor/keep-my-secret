package web

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/config"
	"github.com/grafviktor/keep-my-secret/internal/constant"
	"github.com/grafviktor/keep-my-secret/internal/model"
)

func TestParseMultiPartSecretRequest(t *testing.T) {
	// Create a sample multipart form request with JSON data and a file
	jsonData := `{"name": "mySecret", "description": "Test secret"}`

	// Create a buffer to hold the form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the JSON data as a form field
	_ = writer.WriteField("data", jsonData)

	// Add a sample file to the form
	fileContents := []byte("This is the file content")
	fileWriter, _ := writer.CreateFormFile("file", "sample.txt")
	//nolint:errcheck
	fileWriter.Write(fileContents)

	// Close the form writer
	writer.Close()

	// Create a sample HTTP request with the multipart form data
	req := httptest.NewRequest("POST", "/your-api-endpoint", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a model.Secret instance for storing the parsed data
	secret := &model.Secret{}

	// Call the parseMultiPartSecretRequest function
	err := parseMultiPartSecretRequest(req, secret)
	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify that the parsed data matches the expected values
	expectedSecret := &model.Secret{
		File: fileContents,
	}

	if !reflect.DeepEqual(secret, expectedSecret) {
		t.Errorf("Expected secret %+v, but got %+v", expectedSecret.File, secret.File)
	}
}

func TestSaveSecretHandler(t *testing.T) {
	// Create a sample AppConfig for testing
	appConfig := config.AppConfig{
		// Initialize your AppConfig fields here
		Secret: "your_secret_key",
	}

	// Create a sample HTTP request with JSON data
	jsonData := `{"title": "mySecret", "note": "Test secret"}`
	req, err := http.NewRequest("POST", "/your-api-endpoint", strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), api.ContextUserLogin, "test"))

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Create an instance of your apiRouteProvider with mock dependencies
	handler := &apiRouteProvider{
		config: appConfig,
		storage: &MockStorage{
			users: make(map[string]*model.User),
		},
		keyCache: &MockKeyCache{},
	}

	// Call the SaveSecretHandler
	handler.SaveSecretHandler(rr, req)

	// Check the response status code
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
	}

	// Check the response body
	var response api.Response
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != constant.APIStatusSuccess {
		t.Errorf("Expected API status 'success', got '%s'", response.Status)
	}
}

func TestListSecretsHandler(t *testing.T) {
	// Create a sample AppConfig for testing
	appConfig := config.AppConfig{
		// Initialize your AppConfig fields here
		Secret: "your_secret_key",
	}

	// Create a sample HTTP request
	req, err := http.NewRequest("GET", "/your-api-endpoint", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Create an instance of your apiRouteProvider with mock dependencies
	handler := &apiRouteProvider{
		config: appConfig,
		storage: &MockStorage{
			users: make(map[string]*model.User),
		},
		keyCache: &MockKeyCache{},
	}

	// Create a context with a mock user login value
	ctx := context.WithValue(context.Background(), api.ContextUserLogin, "validLogin")

	// Call the ListSecretsHandler with the mocked context
	handler.ListSecretsHandler(rr, req.WithContext(ctx))

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Parse and check the response body
	var response api.Response
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != constant.APIStatusSuccess {
		t.Errorf("Expected API status 'success', got '%s'", response.Status)
	}
}
