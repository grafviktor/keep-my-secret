package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafviktor/keep-my-secret/internal/config"
)

type mockHandler struct {
	called bool
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestEnableCORS(t *testing.T) {
	// Create a mock HTTP handler for testing
	mock := &mockHandler{}

	// Create a mock HTTP request
	req := httptest.NewRequest("OPTIONS", "http://example.com", nil)

	// Pretend that I'm a kinda browser. Unit test won't set this header for us
	req.Header.Set("Origin", "http://example.com")

	// Create a response recorder
	res := httptest.NewRecorder()

	// Create an instance of the middleware with DevMode enabled
	middlewareInstance := middleware{config: config.AppConfig{DevMode: true}}

	// Call the EnableCORS middleware with the mock handler
	corsHandler := middlewareInstance.EnableCORS(mock)

	// Process the request using the CORS middleware
	corsHandler.ServeHTTP(res, req)

	// Verify that the CORS headers are set correctly when DevMode is enabled
	expectedHeadersDevMode := map[string]string{
		"Access-Control-Allow-Origin":      "http://example.com",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Expose-Headers":    "*",
		"Access-Control-Allow-Methods":     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		"Access-Control-Allow-Headers":     "Accept, Content-Type, X-CSRF-Token, Authorization",
	}

	for header, expectedValue := range expectedHeadersDevMode {
		actualValue := res.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("Expected %s header to be '%s', but got '%s'", header, expectedValue, actualValue)
		}
	}

	// Reset the recorder for the next test
	res = httptest.NewRecorder()

	// Create an instance of the middleware with DevMode disabled
	middlewareInstance.config.DevMode = false

	// Call the EnableCORS middleware with the mock handler
	corsHandler = middlewareInstance.EnableCORS(mock)

	// Process the request using the CORS middleware
	corsHandler.ServeHTTP(res, req)

	// Verify that the CORS headers are not set when DevMode is disabled
	for header := range expectedHeadersDevMode {
		actualValue := res.Header().Get(header)
		if actualValue != "" {
			t.Errorf("Expected %s header to be empty, but got '%s'", header, actualValue)
		}
	}
}
