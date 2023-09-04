package web

import (
	"context"
	"errors"
	"strings"

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
	return nil, nil
}

type mockEncryptor struct{}

func (ms mockEncryptor) Encrypt(plaindata []byte, key string) ([]byte, error) {
	return plaindata, nil
}

func (ms mockEncryptor) Decrypt(cipherdata []byte, key string) ([]byte, error) {
	return cipherdata, nil
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
	// TODO implement me
	panic("implement me")
}

func (mockStorage MockStorage) GetSecret(ctx context.Context, secretID, login string) (*model.Secret, error) {
	// TODO implement me
	panic("implement me")
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
