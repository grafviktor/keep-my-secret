package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/api/auth"
	"github.com/grafviktor/keep-my-secret/internal/config"
)

type mockAuthVerifier struct{}

//nolint:lll
func (m mockAuthVerifier) VerifyAuthHeader(config config.AppConfig, w http.ResponseWriter, r *http.Request) (string, *auth.Claims, error) {
	claims := &auth.Claims{}
	//nolint:goconst
	claims.Subject = "testuser"

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	return "", claims, nil
}

func TestAuthRequired(t *testing.T) {
	// Create an instance of your middleware with the mock dependency
	mw := middleware{
		config:       config.AppConfig{},
		authVerifier: mockAuthVerifier{},
	}

	handler := mw.AuthRequired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user login from the context
		userLogin, ok := r.Context().Value(api.ContextUserLogin).(string)
		if !ok {
			t.Fatal("User login not found in context")
		}

		// Check if the user login matches the expected value
		//noling:goconst
		expectedUserLogin := "testuser" // Replace with the expected user login
		if userLogin != expectedUserLogin {
			t.Errorf("Expected user login '%s', got '%s'", expectedUserLogin, userLogin)
		}

		// Serve the response
		w.WriteHeader(http.StatusOK)
	}))

	// Case 1 verifier should pass is OK
	// Create a sample HTTP request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer bla-bla")

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Call the AuthRequired middleware
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Case 2 verifier should send an error
	// Create a sample HTTP request
	req, err = http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Not setting auth header this time
	// req.Header.Set("Authorization", "Bearer bla-bla")

	// Create a mock response recorder
	rr = httptest.NewRecorder()

	// Call the AuthRequired middleware
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}
