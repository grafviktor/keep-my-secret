package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grafviktor/keep-my-secret/internal/api/auth"

	"github.com/grafviktor/keep-my-secret/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafviktor/keep-my-secret/internal/config"
	"github.com/grafviktor/keep-my-secret/internal/storage"
)

var appConfig = config.AppConfig{
	StorageType: storage.TypeSQL,
}

func TestUserHTTPHandler_Register(t *testing.T) {
	type expected struct {
		statusCode  int
		contentType string
		headerKey   string
		headerValue string
		// response    string
	}

	type httpResponseTestCase struct {
		name string
		// httpPath    string
		httpBody    string
		httpMethod  string
		httpHandler http.HandlerFunc
		expected    expected
		// headerKey   string
		// headerValue string
	}

	storage := &MockStorage{
		users: make(map[string]*model.User),
	}
	handler := newUserHandlerProvider(appConfig, storage)
	urlPath := "/api/v1/user/register"

	SuccessfulLogin := httpResponseTestCase{
		name:        "Register a new credentials success",
		httpHandler: handler.RegisterHandler,
		httpMethod:  http.MethodPost,
		httpBody:    `{"username":"tony.tester@example.com", "password":"1"}`,
		expected: expected{
			contentType: "application/json",
			statusCode:  http.StatusCreated,
		},
	}

	ErrorDuplicateRecord := httpResponseTestCase{
		name:        "Register the credentials already exists error",
		httpHandler: handler.RegisterHandler,
		httpMethod:  http.MethodPost,
		httpBody:    `{"username":"tony.tester@example.com", "password":"1"}`,
		expected: expected{
			contentType: "application/json",
			statusCode:  http.StatusConflict,
		},
	}

	ErrorNoUsernameProvided := httpResponseTestCase{
		name:        "Register no username provided error",
		httpHandler: handler.RegisterHandler,
		httpMethod:  http.MethodPost,
		httpBody:    `{"password":"1"}`,
		expected: expected{
			contentType: "application/json",
			statusCode:  http.StatusNotAcceptable,
		},
	}

	ErrorNoPasswordProvided := httpResponseTestCase{
		name:        "Register no password provided error",
		httpHandler: handler.RegisterHandler,
		httpMethod:  http.MethodPost,
		httpBody:    `{"username":"tony.tester@example.com", "password":""}`,
		expected: expected{
			contentType: "application/json",
			statusCode:  http.StatusNotAcceptable,
		},
	}

	testCases := []httpResponseTestCase{
		SuccessfulLogin,
		ErrorDuplicateRecord,
		ErrorNoUsernameProvided,
		ErrorNoPasswordProvided,
	}

	r := NewHTTPRouter(appConfig, storage)
	ts := httptest.NewServer(r)
	defer ts.Close()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			body := strings.NewReader(testCase.httpBody)
			request, err := http.NewRequest(testCase.httpMethod, ts.URL+urlPath, body)
			require.NoError(t, err)

			received, err := client.Do(request)
			require.NoError(t, err)

			assert.Equal(t, testCase.expected.statusCode, received.StatusCode)
			assert.Equal(t, testCase.expected.headerValue, received.Header.Get(testCase.expected.headerKey))
			assert.Equal(t, testCase.expected.contentType, received.Header.Get("Content-Type"))

			_ = received.Body.Close()
		})
	}
}

var _ storage.Storage = &MockStorage{}

func TestUserHTTPHandler_Login(t *testing.T) {
	type expected struct {
		statusCode  int
		contentType string
		headerKey   string
		headerValue string
		// response    string
	}

	type httpResponseTestCase struct {
		name string
		// httpPath    string
		httpBody    string
		httpMethod  string
		httpHandler http.HandlerFunc
		expected    expected
		// headerKey   string
		// headerValue string
	}

	ls := MockStorage{
		users: make(map[string]*model.User),
	}

	ls.users["tony.tester@example.com"] = &model.User{
		Login:           "tony.tester@example.com",
		HashedPassword:  "$2a$10$AokZyUVIqfgBtEwCNhOzbeE68Zk6uwZ42NvDdPK24Xesmb08OJ.DO",
		RestorePassword: "",
	}

	handler := newUserHandlerProvider(appConfig, &ls)
	urlPath := "/api/v1/user/login"

	// SuccessfulLogin := httpResponseTestCase{
	//	name:        "Login credentials success",
	//	httpHandler: handler.LoginHandler,
	//	httpMethod:  http.MethodPost,
	//	httpBody:    `{"username":"tony.tester@example.com", "password":"1"}`,
	//	expected: expected{
	//		contentType: "application/json",
	//		statusCode:  http.StatusOK,
	//	},
	// }

	UnsuccessfulLogin := httpResponseTestCase{
		name:        "Login credentials error",
		httpHandler: handler.LoginHandler,
		httpMethod:  http.MethodPost,
		httpBody:    `{"username":"tony.tester@example.com", "password":"12"}`,
		expected: expected{
			contentType: "application/json",
			statusCode:  http.StatusUnauthorized,
		},
	}

	MissingAttributesRequestLogin := httpResponseTestCase{
		name:        "Login credentials required attributes missing error",
		httpHandler: handler.LoginHandler,
		httpMethod:  http.MethodPost,
		httpBody:    `{"username":"tony.tester@example.com"}`,
		expected: expected{
			contentType: "application/json",
			statusCode:  http.StatusUnauthorized,
		},
	}

	BadRequestLogin := httpResponseTestCase{
		name:        "Login credentials malformed body error",
		httpHandler: handler.LoginHandler,
		httpMethod:  http.MethodPost,
		httpBody:    "",
		expected: expected{
			contentType: "application/json",
			statusCode:  http.StatusBadRequest,
		},
	}

	testCases := []httpResponseTestCase{
		// SuccessfulLogin,
		UnsuccessfulLogin,
		MissingAttributesRequestLogin,
		BadRequestLogin,
	}

	r := NewHTTPRouter(appConfig, &ls)
	ts := httptest.NewServer(r)
	defer ts.Close()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			body := strings.NewReader(testCase.httpBody)
			request, err := http.NewRequest(testCase.httpMethod, ts.URL+urlPath, body)
			require.NoError(t, err)

			received, err := client.Do(request)
			require.NoError(t, err)

			assert.Equal(t, testCase.expected.statusCode, received.StatusCode)
			assert.Equal(t, testCase.expected.headerValue, received.Header.Get(testCase.expected.headerKey))
			assert.Equal(t, testCase.expected.contentType, received.Header.Get("Content-Type"))

			_ = received.Body.Close()
		})
	}
}

type MockUser struct{}

func (u *MockUser) GetDataKey(password string) (string, error) {
	//nolint:goconst
	return "mocked-secret", nil
}

type MockKeyCache struct {
	setCalled      bool
	setLogin       string
	setSecret      string
	getCalled      bool
	getLogin       string
	getReturnValue string
}

func (kc *MockKeyCache) Set(login, secret string) {
	kc.setCalled = true
	kc.setLogin = login
	kc.setSecret = secret
}

func (kc *MockKeyCache) Get(login string) (string, error) {
	kc.getCalled = true
	kc.getLogin = login
	return kc.getReturnValue, nil
}

type MockAuthUtils struct {
	generateTokenPairCalled bool
	generateTokenPairUser   *auth.JWTUser
	generateTokenPairReturn auth.TokenPair
	getRefreshCookieCalled  bool
	getRefreshCookieToken   string
	getRefreshCookieReturn  *http.Cookie
	shouldTriggerError      bool
}

func (au *MockAuthUtils) GenerateTokenPair(user *auth.JWTUser) (auth.TokenPair, error) {
	if au.shouldTriggerError {
		return au.generateTokenPairReturn, errors.New("that is a mock error triggered by 'error' subject in token")
	}

	au.generateTokenPairCalled = true
	au.generateTokenPairUser = user

	return au.generateTokenPairReturn, nil
}

func (au *MockAuthUtils) GetRefreshCookie(token string) *http.Cookie {
	au.getRefreshCookieCalled = true
	au.getRefreshCookieToken = token
	return au.getRefreshCookieReturn
}

func TestHandleSuccessFullUserSignIn(t *testing.T) {
	handler := &userHTTPHandler{
		config:    config.AppConfig{},
		keyCache:  &MockKeyCache{},
		authUtils: &MockAuthUtils{},
	}

	httptest.NewRequest("POST", "/signin", nil)
	w := httptest.NewRecorder()

	cred := credentials{
		Login:    "testuser",
		Password: "testpassword",
	}

	handler.handleSuccessFullUserSignIn(w, &MockUser{}, cred)

	// Assert the behavior of MockKeyCache and MockAuthUtils
	if !handler.keyCache.(*MockKeyCache).setCalled {
		t.Error("Expected Set to be called on keyCache")
	}
	//nolint:goconst
	if handler.keyCache.(*MockKeyCache).setLogin != "testuser" {
		t.Errorf("Expected Set to be called with login 'testuser', but got '%s'", handler.keyCache.(*MockKeyCache).setLogin)
	}
	if handler.keyCache.(*MockKeyCache).setSecret != "mocked-secret" {
		//nolint:lll
		t.Errorf("Expected Set to be called with secret 'mocked-secret', but got '%s'", handler.keyCache.(*MockKeyCache).setSecret)
	}

	if !handler.authUtils.(*MockAuthUtils).generateTokenPairCalled {
		t.Error("Expected GenerateTokenPair to be called on authUtils")
	}
	if handler.authUtils.(*MockAuthUtils).generateTokenPairUser.ID != "testuser" {
		//nolint:lll
		t.Errorf("Expected GenerateTokenPair to be called with user ID 'testuser', but got '%s'", handler.authUtils.(*MockAuthUtils).generateTokenPairUser.ID)
	}

	// Assert response status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}
}

// Create a mockTokenGenerator that implements the GenerateTokenPair method
// type mockTokenGenerator struct{}
//
// func (a *mockTokenGenerator) GenerateTokenPair(user *auth.JWTUser) (*auth.TokenPair, error) {
// 	// Implement your mock logic here
// 	// For example, you can return a mock token pair or an error based on test cases
// 	if user.ID == "validUserID" {
// 		return &auth.TokenPair{
// 			AccessToken:  "mockAccessToken",
// 			RefreshToken: "mockRefreshToken",
// 		}, nil
// 	}
// 	return nil, errors.New("mockTokenGenerator: error")
// }

func TestRefreshTokenHandler(t *testing.T) {
	validTokenSubject := "1234567890" // that's taken from token (decrypt the token on jwt.io, and you'll see this subject)
	//nolint:lll
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	testCases := []struct {
		name           string
		secret         string
		token          string
		responseStatus int
		tokenSubject   string
		authUtilsError bool
	}{
		{
			name:           "Valid token",
			secret:         "your-256-bit-secret", // Copied from jwt.io
			token:          validToken,
			responseStatus: http.StatusOK,
			tokenSubject:   validTokenSubject,
		}, {
			name:   "Invalid token",
			secret: "your-256-bit-secret",
			//nolint:lll
			token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			responseStatus: http.StatusUnauthorized,
			tokenSubject:   validTokenSubject,
		}, {
			name:           "Username not found in storage",
			secret:         "your-256-bit-secret",
			token:          validToken,
			responseStatus: http.StatusUnauthorized,
			tokenSubject:   "the_username_id_different_from_what_token_claims",
		}, {
			name:           "AuthUtils triggered error",
			secret:         "your-256-bit-secret",
			token:          validToken,
			responseStatus: http.StatusInternalServerError,
			tokenSubject:   validTokenSubject,
			authUtilsError: true, // will trigger error in auth utils
		},
	}

	for _, testCase := range testCases {
		// Create a sample AppConfig for testing
		appConfig := config.AppConfig{
			// Initialize your AppConfig fields here
			Secret: testCase.secret,
		}

		// Create a sample HTTP request with a cookie
		req, err := http.NewRequest("GET", "/refresh", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: testCase.token})

		// Create a mock response recorder
		rr := httptest.NewRecorder()

		// Create an instance of your userHTTPHandler with the mock dependencies
		handler := &userHTTPHandler{
			config: appConfig,
			storage: &MockStorage{
				users: make(map[string]*model.User),
			},
			authUtils: &MockAuthUtils{
				shouldTriggerError: testCase.authUtilsError,
			},
		}

		//nolint:errcheck
		handler.storage.AddUser(req.Context(), &model.User{
			Login: testCase.tokenSubject,
		})

		// Call the RefreshTokenHandler
		handler.RefreshTokenHandler(rr, req)

		// Check the response status code
		if rr.Code != testCase.responseStatus {
			t.Errorf("%s: Expected status code %d, got %d", testCase.name, testCase.responseStatus, rr.Code)
		}
	}
}

func TestLogoutHandler(t *testing.T) {
	// Create a sample AppConfig for testing
	appConfig := config.AppConfig{
		// Initialize your AppConfig fields here
		Secret: "your_secret_key",
	}

	// Create a sample HTTP request
	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Create an instance of your userHTTPHandler with the mock dependency
	handler := &userHTTPHandler{
		config: appConfig,
	}

	// Call the LogoutHandler
	handler.LogoutHandler(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Check the cookie in the response
	cookies := rr.Result().Cookies()
	if len(cookies) != 1 {
		t.Error("Expected one cookie in the response")
	} else {
		// Check if the cookie is expired
		if cookies[0].MaxAge != -1 {
			t.Error("Expected the cookie to be expired")
		}
		// Check if the cookie name matches
		if cookies[0].Name != auth.CookieName {
			t.Errorf("Expected cookie name to be %s, got %s", auth.CookieName, cookies[0].Name)
		}
		// Check if the cookie value matches
		// if cookies[0].Value != "expiredRefreshToken" {
		// 	t.Errorf("Expected cookie value to be 'expiredRefreshToken', got '%s'", cookies[0].Value)
		// }
	}
}
