package storage

import (
	"context"
	"errors"

	"github.com/grafviktor/keep-my-secret/internal/model"
	"github.com/grafviktor/keep-my-secret/internal/storage/sql"
)

type Storage interface {
	AddUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, login string) (*model.User, error)
	SaveSecret(ctx context.Context, secret *model.Secret, login string) (*model.Secret, error)
	GetSecretsByUser(ctx context.Context, login string) (map[int]*model.Secret, error)
	DeleteSecret(ctx context.Context, secretID, login string) error
	GetSecret(ctx context.Context, secretID, login string) (*model.Secret, error)
	Close() error
}

func GetStorage(ctx context.Context, storageType Type, dsn string) (Storage, error) {
	if storageType != TypeSQL {
		return nil, errors.New("unsupported storage type")
	}

	return sql.NewSQLStorage(ctx, dsn), nil
}
