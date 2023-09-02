package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"

	"github.com/grafviktor/keep-my-secret/internal/config"
)

const cookieName = "refresh_cookie"

func TestGetRefreshCookie(t *testing.T) {
	// Initialize your Auth object with relevant settings
	auth := &Auth{
		CookieName:    cookieName,
		CookiePath:    "/",
		RefreshExpiry: 24 * time.Hour,
	}

	// Call the GetRefreshCookie function
	refreshToken := "your_refresh_token"
	cookie := auth.GetRefreshCookie(refreshToken)

	// Assert that the cookie has the expected values
	if cookie.Name != cookieName {
		t.Errorf("Expected cookie name '%s', got '%s'", cookieName, cookie.Name)
	}

	if cookie.Path != "/" {
		t.Errorf("Expected cookie path '/', got '%s'", cookie.Path)
	}

	if cookie.Value != refreshToken {
		t.Errorf("Expected cookie value '%s', got '%s'", refreshToken, cookie.Value)
	}

	if cookie.MaxAge != int(24*60*60) {
		t.Errorf("Expected cookie MaxAge '%d', got '%d'", int(24*60*60), cookie.MaxAge)
	}

	if cookie.SameSite != siteMode {
		t.Errorf("Expected SameSite mode '%v', got '%v'", siteMode, cookie.SameSite)
	}

	if !cookie.HttpOnly {
		t.Error("Expected HttpOnly to be true")
	}

	if !cookie.Secure {
		t.Error("Expected Secure to be true")
	}
}

func TestGetExpiredRefreshCookie(t *testing.T) {
	// Initialize your Auth object with relevant settings
	auth := &Auth{
		CookieName: "refresh_cookie",
		CookiePath: "/",
	}

	// Call the GetExpiredRefreshCookie function
	cookie := auth.GetExpiredRefreshCookie()

	// Assert that the cookie has the expected values
	if cookie.Name != "refresh_cookie" {
		t.Errorf("Expected cookie name 'refresh_cookie', got '%s'", cookie.Name)
	}

	if cookie.Path != "/" {
		t.Errorf("Expected cookie path '/', got '%s'", cookie.Path)
	}

	if cookie.Value != "" {
		t.Errorf("Expected empty cookie value, got '%s'", cookie.Value)
	}

	if !cookie.Expires.Equal(time.Unix(0, 0)) {
		t.Errorf("Expected cookie expiry '1970-01-01 00:00:00 UTC', got '%s'", cookie.Expires)
	}

	if cookie.MaxAge != -1 {
		t.Errorf("Expected MaxAge '-1', got '%d'", cookie.MaxAge)
	}

	if cookie.SameSite != siteMode {
		t.Errorf("Expected SameSite mode '%v', got '%v'", siteMode, cookie.SameSite)
	}

	if !cookie.HttpOnly {
		t.Error("Expected HttpOnly to be true")
	}

	if !cookie.Secure {
		t.Error("Expected Secure to be true")
	}
}

func TestVerifyAuthHeader(t *testing.T) {
	envConfig := config.EnvConfig{
		Secret: "romeo romeo whiskey",
		Domain: "localhost",
	}
	ac := config.New(envConfig)

	// JWTSecret := ac.Secret
	// JWTIssuer := ac.JWTIssuer

	// Create a test HTTP request with an Authorization header
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	jwtUser := JWTUser{ID: "user@localhost"}
	newAuth := New(ac)
	tokenPair, err := newAuth.GenerateTokenPair(&jwtUser)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	tokenString := tokenPair.AccessToken
	// tokenString := GenerateTokenPair(Secret, JWTIssuer, time.Now().Add(1*time.Hour))
	req.Header.Add("Authorization", "Bearer "+tokenString)

	// Create a test ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the VerifyAuthHeader function
	token, claims, err := VerifyAuthHeader(ac, rr, req)
	// Check for expected results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if token != tokenString {
		t.Errorf("Expected token '%s', got '%s'", tokenString, token)
	}

	require.NotNil(t, claims)

	// if claims == nil {
	// 	t.Error("Expected non-nil claims")
	// }

	if claims.Issuer != ac.JWTIssuer {
		t.Errorf("Expected issuer '%s', got '%s'", ac.JWTIssuer, claims.Issuer)
	}

	// Test cases with invalid headers
	testCases := []struct {
		headerValue string
		expectedErr string
	}{
		{"", "no auth header"},
		{"InvalidHeader", "invalid auth header"},
		{"Bearer InvalidToken", "token contains an invalid number of segments"},
		{"Bearer " + generateJWTToken("WrongSecret", ac.JWTIssuer, time.Now().Add(1*time.Hour)), "signature is invalid"},
		{"Bearer " + generateJWTToken(ac.Secret, ac.JWTIssuer, time.Now().Add(-1*time.Hour)), "expired token"},
		{"Bearer " + generateJWTToken(ac.Secret, "WrongIssuer", time.Now().Add(1*time.Hour)), "invalid issuer"},
	}

	for _, tc := range testCases {
		req.Header.Set("Authorization", tc.headerValue)
		rr = httptest.NewRecorder()
		_, _, err := VerifyAuthHeader(ac, rr, req)
		if err == nil || err.Error() != tc.expectedErr {
			t.Errorf("Expected error: '%s', got: '%v'", tc.expectedErr, err)
		}
	}
}

func generateJWTToken(secretKey, issuer string, expiration time.Time) string {
	claims := jwt.MapClaims{
		"iss": issuer,
		"exp": expiration.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}
