package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grafviktor/keep-my-secret/internal/constant"
	"github.com/grafviktor/keep-my-secret/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafviktor/keep-my-secret/internal/config"
	"github.com/grafviktor/keep-my-secret/internal/storage"
)

var appConfig = config.AppConfig{
	StorageType: storage.TypeSQL,
}

var testContext = context.Background()

var _ storage.Storage = MockStorage{}

type MockStorage struct {
	users map[string]*model.User
}

func (mockStorage MockStorage) SaveSecret(ctx context.Context, secret *model.Secret, login string) (*model.Secret, error) {
	//TODO implement me
	panic("implement me")
}

func (mockStorage MockStorage) GetSecretsByUser(ctx context.Context, login string) (map[int]*model.Secret, error) {
	//TODO implement me
	panic("implement me")
}

func (mockStorage MockStorage) DeleteSecret(ctx context.Context, secretID, login string) error {
	//TODO implement me
	panic("implement me")
}

func (mockStorage MockStorage) GetSecret(ctx context.Context, secretID, login string) (*model.Secret, error) {
	//TODO implement me
	panic("implement me")
}

func (mockStorage MockStorage) Close() error {
	//TODO implement me
	panic("implement me")
}

func (mockStorage MockStorage) AddUser(ctx context.Context, user *model.User) (*model.User, error) {
	if strings.Trim(user.Login, " ") == "" {
		return nil, constant.ErrBadArgument
	}

	_, ok := mockStorage.users[user.Login]
	if ok {
		return nil, constant.ErrDuplicateRecord
	}

	mockStorage.users[user.Login] = user

	return user, nil
}

func (mockStorage MockStorage) GetUser(ctx context.Context, login string) (*model.User, error) {
	//if login == "tony.tester@example.com" {
	//	return &model.User{
	//		Login:           "tony.tester@example.com",
	//		HashedPassword:  "1",
	//		RestorePassword: "",
	//	}, nil
	//}

	user, ok := mockStorage.users[login]
	if ok {
		return user, nil
	}

	return nil, constant.ErrNotFound
}

func TestUserHTTPHandler_Register(t *testing.T) {
	type expected struct {
		statusCode  int
		contentType string
		headerKey   string
		headerValue string
		response    string
	}

	type httpResponseTestCase struct {
		name        string
		httpPath    string
		httpBody    string
		httpMethod  string
		httpHandler http.HandlerFunc
		expected    expected
		headerKey   string
		headerValue string
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
		response    string
	}

	type httpResponseTestCase struct {
		name        string
		httpPath    string
		httpBody    string
		httpMethod  string
		httpHandler http.HandlerFunc
		expected    expected
		headerKey   string
		headerValue string
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

	//SuccessfulLogin := httpResponseTestCase{
	//	name:        "Login credentials success",
	//	httpHandler: handler.LoginHandler,
	//	httpMethod:  http.MethodPost,
	//	httpBody:    `{"username":"tony.tester@example.com", "password":"1"}`,
	//	expected: expected{
	//		contentType: "application/json",
	//		statusCode:  http.StatusOK,
	//	},
	//}

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
		//SuccessfulLogin,
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
