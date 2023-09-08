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

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestParseMultiPartSecretRequestNegative(t *testing.T) {
	// Malformed JSON
	jsonData := `{"name": "mySecret", faulty`

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
	require.Error(t, err)

	// Wrong content type
	req = httptest.NewRequest("POST", "/your-api-endpoint", body)
	req.Header.Set("Content-Type", "")
	err = parseMultiPartSecretRequest(req, nil)
	require.Error(t, err)
}

func TestSaveSecretHandler(t *testing.T) {
	// Create a sample AppConfig for testing
	appConfig := config.AppConfig{
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

	// Create an instance of my apiRouteProvider with mock dependencies
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

func TestSaveSecretHandlerNegative(t *testing.T) {
	// Create a sample AppConfig for testing
	appConfig := config.AppConfig{
		Secret: "your_secret_key",
	}

	testCases := []struct {
		name                      string
		payload                   string
		login                     string
		shouldSetContextUserLogin bool
		httpStatusCode            int
	}{
		{
			name:                      "no user login in request context",
			payload:                   `{"title": "mySecret", "note": "Test secret"}`,
			shouldSetContextUserLogin: false,
			login:                     "valid_user",
			httpStatusCode:            http.StatusUnauthorized,
		},
		{
			name:                      "malformed JSON",
			payload:                   `{"title": "mySecret", malformed`,
			shouldSetContextUserLogin: true,
			login:                     "valid_user",
			httpStatusCode:            http.StatusBadRequest,
		},
		{
			name:                      "error getting decrypt key from keycache",
			payload:                   `{"title": "mySecret", "note": "Test secret"}`,
			shouldSetContextUserLogin: true,
			login:                     "invalid_user",
			httpStatusCode:            http.StatusUnauthorized,
		},
		{
			name:                      "no user login in request context",
			payload:                   `{"title": "mySecret", "note": "Test secret"}`,
			shouldSetContextUserLogin: true,
			login:                     "valid_user_invalid_secret",
			httpStatusCode:            http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		// Create a sample HTTP request with JSON data
		req, _ := http.NewRequest("POST", "/your-api-endpoint", strings.NewReader(tc.payload))
		req.Header.Set("Content-Type", "application/json")

		if tc.shouldSetContextUserLogin {
			req = req.WithContext(context.WithValue(req.Context(), api.ContextUserLogin, tc.login))
		}

		// Create a mock response recorder
		rr := httptest.NewRecorder()

		// Create an instance of my apiRouteProvider with mock dependencies
		handler := &apiRouteProvider{
			config: appConfig,
			storage: &MockStorage{
				users: make(map[string]*model.User),
			},
			keyCache: &MockKeyCache{},
		}

		handler.SaveSecretHandler(rr, req)

		require.Equal(t, tc.httpStatusCode, rr.Code)
	}
}

func TestListSecretsHandler(t *testing.T) {
	// Create a sample AppConfig for testing
	appConfig := config.AppConfig{
		Secret: "your_secret_key",
	}

	// Create a sample HTTP request
	req, err := http.NewRequest("GET", "/your-api-endpoint", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Create an instance of my apiRouteProvider with mock dependencies
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

func TestListSecretsHandlerNegative(t *testing.T) {
	appConfig := config.AppConfig{
		Secret: "your_secret_key",
	}

	req, err := http.NewRequest("GET", "/api-endpoint", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	handler := &apiRouteProvider{
		config: appConfig,
		storage: &MockStorage{
			users: make(map[string]*model.User),
		},
		keyCache: &MockKeyCache{},
	}

	rr := httptest.NewRecorder()

	ctx := context.WithValue(context.Background(), api.ContextUserLogin, "invalid_user")
	handler.ListSecretsHandler(rr, req.WithContext(ctx))
	require.Equal(t, rr.Code, http.StatusUnauthorized)

	rr = httptest.NewRecorder()

	ctx = context.WithValue(context.Background(), api.ContextUserLogin, "valid_user")
	handler.ListSecretsHandler(rr, req.WithContext(ctx))
	require.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestDeleteSecretHandler(t *testing.T) {
	// Create a new chi router for testing.
	r := chi.NewRouter()

	// Create a new mock storage instance.
	mock := &MockStorage{}

	// Create an instance of the API route provider with the mock storage.
	apiProvider := &apiRouteProvider{storage: mock}

	// Register the DeleteSecretHandler route on the router.
	r.Delete("/secrets/{id}", apiProvider.DeleteSecretHandler)

	// Create a test request with a valid ID.
	req := httptest.NewRequest("DELETE", "/secrets/valid_id", nil)
	ctx := context.WithValue(req.Context(), api.ContextUserLogin, "test_user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute the request.
	r.ServeHTTP(w, req)

	// Check the response status code.
	assert.Equal(t, http.StatusAccepted, w.Code)

	// Check the response body.
	expectedResponse := `{"status":"success","message":"","data":"valid_id"}`
	assert.Equal(t, expectedResponse, w.Body.String())

	// Create a test request with an invalid ID.
	req = httptest.NewRequest("DELETE", "/secrets/invalid_id", nil)
	ctx = context.WithValue(req.Context(), api.ContextUserLogin, "test_user")
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()

	// Execute the request.
	r.ServeHTTP(w, req)

	// Check the response status code.
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Check the response body.
	expectedResponse = `{"status":"fail","message":"not found","data":null}`
	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestDownloadSecretFileHandler(t *testing.T) {
	// Create a new chi router for testing.
	r := chi.NewRouter()

	// Create a new mock storage instance and a mock key cache instance.
	mockStorage := &MockStorage{}
	mockKeyCache := &MockKeyCache{}

	// Create an instance of the API route provider with the mock dependencies.
	apiProvider := &apiRouteProvider{
		storage:  mockStorage,
		keyCache: mockKeyCache,
	}

	// Register the DownloadSecretFileHandler route on the router.
	r.Get("/secrets/{id}", apiProvider.DownloadSecretFileHandler)

	// Create a test request with a valid ID and valid user.
	req := httptest.NewRequest("GET", "/secrets/valid_id", nil)
	ctx := context.WithValue(req.Context(), api.ContextUserLogin, "valid_user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute the request.
	r.ServeHTTP(w, req)

	// Check the response status code.
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response headers.
	assert.Equal(t, "attachment; filename=test.txt", w.Header().Get("Content-Disposition"))
	assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, "20", w.Header().Get("Content-Length"))

	// Check the response body.
	expectedResponse := "This is a test file."
	assert.Equal(t, expectedResponse, w.Body.String())

	// Create a test request with an invalid ID and valid user.
	req = httptest.NewRequest("GET", "/secrets/not_found_id", nil)
	ctx = context.WithValue(req.Context(), api.ContextUserLogin, "valid_user")
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()

	// Execute the request.
	r.ServeHTTP(w, req)

	// Check the response status code for not found.
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Create a test request with a valid ID and invalid user.
	req = httptest.NewRequest("GET", "/secrets/valid_id", nil)
	ctx = context.WithValue(req.Context(), api.ContextUserLogin, "invalid_user")
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()

	// Execute the request.
	r.ServeHTTP(w, req)

	// Check the response status code for unauthorized.
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Create a test request with an error scenario (mock storage error).
	req = httptest.NewRequest("GET", "/secrets/valid_id", nil)
	// ctx = context.WithValue(req.Context(), api.ContextUserLogin, nil)
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()

	// Execute the request.
	r.ServeHTTP(w, req)

	// Check the response status code for internal server error.
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
