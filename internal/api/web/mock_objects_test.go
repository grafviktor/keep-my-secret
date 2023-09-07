package web

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/grafviktor/keep-my-secret/internal/api/auth"

	"github.com/grafviktor/keep-my-secret/internal/constant"
	"github.com/grafviktor/keep-my-secret/internal/model"
	"github.com/grafviktor/keep-my-secret/internal/storage"
)

var (
	_ storage.Storage = MockStorage{}
	_ model.Encryptor = mockEncryptor{}
)

type MockStorage struct {
	users map[string]*model.User
}

//nolint:lll
func (mockStorage MockStorage) SaveSecret(ctx context.Context, secret *model.Secret, login string) (*model.Secret, error) {
	if login == "valid_user_invalid_secret" {
		return nil, errors.New("mock storage error for invalid secret")
	}
	return nil, nil
}

type mockEncryptor struct{}

func (ms mockEncryptor) Encrypt(secret *model.Secret, key, salt string) error {
	return nil
}

func (ms mockEncryptor) Decrypt(secret *model.Secret, key, salt string) error {
	return nil
}

func (mockStorage MockStorage) GetSecretsByUser(ctx context.Context, login string) (map[int]*model.Secret, error) {
	if login == "validLogin" {
		// Create and return mock secrets
		secret1 := &model.Secret{
			ID:             0,
			Type:           "",
			Title:          "",
			Login:          "",
			Password:       "",
			Note:           "",
			File:           []byte("Mock file content 1"),
			FileName:       "",
			CardholderName: "",
			CardNumber:     "",
			Expiration:     "",
			SecurityCode:   "",
		}
		secret2 := &model.Secret{
			ID:             0,
			Type:           "",
			Title:          "",
			Login:          "",
			Password:       "",
			Note:           "",
			File:           []byte("Mock file content 2"),
			FileName:       "",
			CardholderName: "",
			CardNumber:     "",
			Expiration:     "",
			SecurityCode:   "",
		}
		secret1.SetEncryptor(mockEncryptor{})
		secret2.SetEncryptor(mockEncryptor{})
		return map[int]*model.Secret{
			1: secret1,
			2: secret2,
		}, nil
	}
	return nil, errors.New("mockStorage: error")
}

func (mockStorage MockStorage) DeleteSecret(ctx context.Context, secretID, login string) error {
	if secretID == "invalid_id" {
		return constant.ErrNotFound
	}

	return nil
}

func (mockStorage MockStorage) GetSecret(ctx context.Context, secretID, login string) (*model.Secret, error) {
	// Simulate fetching a secret based on the test scenario.
	//nolint:gocritic
	if secretID == "valid_id" {
		secret := &model.Secret{
			ID:             0,
			Type:           "file",
			Title:          "Test",
			Login:          "tony@tester",
			Password:       "",
			Note:           "",
			File:           []byte("This is a test file."),
			FileName:       "test.txt",
			CardholderName: "",
			CardNumber:     "",
			Expiration:     "",
			SecurityCode:   "",
			Encryptor:      nil,
		}

		secret.SetEncryptor(mockEncryptor{})
		return secret, nil
	} else if secretID == "not_found_id" {
		return nil, constant.ErrNotFound
	} else {
		return nil, errors.New("mock storage error")
	}
}

func (mockStorage MockStorage) Close() error {
	// TODO implement me
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
	user, ok := mockStorage.users[login]
	if ok {
		return user, nil
	}

	return nil, constant.ErrNotFound
}

type MockUser struct{}

func (u *MockUser) GetDataKey(password string) (string, error) {
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
	if login == "invalid_user" {
		return "", errors.New("mock key cache error")
	}

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
