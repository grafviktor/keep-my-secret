package web

import (
	"context"
	"net/http"

	"github.com/grafviktor/keep-my-secret/internal/api/auth"
	"github.com/grafviktor/keep-my-secret/internal/model"
)

type userStorage interface {
	AddUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, login string) (*model.User, error)
}

type keyCache interface {
	Set(login, key string)
	Get(login string) (string, error)
}

type authUtils interface {
	GenerateTokenPair(user *auth.JWTUser) (auth.TokenPair, error)
	GetRefreshCookie(token string) *http.Cookie
}
