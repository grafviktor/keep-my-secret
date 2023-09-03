package middleware

/*
import (
	"context"
	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/api/auth"
	"github.com/grafviktor/keep-my-secret/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockAuth struct {
}

func (m mockAuth) AuthRequired(fn http.HandlerFunc) {
	// Mock authentication logic here (return a valid Claims struct for success)
	fn()
}

// Mock implementation of the auth.VerifyAuthHeader function
func (m mockAuth) mockVerifyAuthHeader(_ config.AppConfig, w http.ResponseWriter, r *http.Request)
(*auth.Claims, error) {
	// Mock authentication logic here (return a valid Claims struct for success)
	claims := auth.Claims{}
	claims.Subject = "tony.tester@locahost"
	return &claims, nil
}

type mockMiddleware struct {
	config   config.AppConfig
	mockAuth mockAuth
	middleware
}

func TestAuthRequired(t *testing.T) {
	// Create a mock HTTP handler for testing
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLogin := r.Context().Value(api.ContextUserLogin).(string)
		w.Write([]byte("Authenticated user: " + userLogin))
	})

	// Create a mock HTTP request
	req := httptest.NewRequest("GET", "http://example.com/some/protected/resource", nil)
	req = req.WithContext(context.WithValue(req.Context(), api.ContextUserLogin, "tony.tester@locahost"))

	// Create a response recorder
	res := httptest.NewRecorder()

	// Create an instance of the middleware with mockVerifyAuthHeader
	middlewareInstance := mockMiddleware{
		config:   config.AppConfig{},
		mockAuth: mockAuth{},
	}
	// middlewareInstance.auth.VerifyAuthHeader = mockVerifyAuthHeader

	// Call the AuthRequired middleware with the mock handler
	authRequiredHandler := middlewareInstance.mockAuth.AuthRequired(mockHandler)

	// Process the request using the AuthRequired middleware
	authRequiredHandler.ServeHTTP(res, req)

	// Verify that the expected user login is set in the request context
	expectedUserLogin := "tony.tester@locahost"
	userLogin := req.Context().Value(api.ContextUserLogin).(string)
	if userLogin != expectedUserLogin {
		t.Errorf("Expected user login to be '%s', but got '%s'", expectedUserLogin, userLogin)
	}

	// Verify the response from the mock handler
	expectedResponse := "Authenticated user: user123"
	if res.Body.String() != expectedResponse {
		t.Errorf("Expected response body to be '%s', but got '%s'", expectedResponse, res.Body.String())
	}

	// Reset the recorder for the next test
	res = httptest.NewRecorder()

	// Mock authentication failure by returning an error from mockVerifyAuthHeader
	// middlewareInstance.auth.VerifyAuthHeader = func(config config.AppConfig, w http.ResponseWriter, r *http.Request)
(*auth.Claims, error) {
	// 	return nil, errors.New(constant.APIMessageUnauthorized)
	// }

	// Call the AuthRequired middleware with the mock handler
	authRequiredHandler = middlewareInstance.AuthRequired(mockHandler)

	// Process the request using the AuthRequired middleware
	authRequiredHandler.ServeHTTP(res, req)

	// Verify that the response indicates unauthorized access
	if res.Code != http.StatusUnauthorized {
		t.Errorf("Expected HTTP status code %d, but got %d", http.StatusUnauthorized, res.Code)
	}
}
*/
