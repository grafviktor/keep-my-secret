package middleware

import (
	"net/http"

	"github.com/grafviktor/keep-my-secret/internal/api/auth"
	"github.com/grafviktor/keep-my-secret/internal/config"
)

type TokenVerifier interface {
	VerifyAuthHeader(config config.AppConfig, w http.ResponseWriter, r *http.Request) (string, *auth.Claims, error)
}

type middleware struct {
	config       config.AppConfig
	authVerifier TokenVerifier
}

func New(appConfig config.AppConfig) middleware {
	return middleware{
		config:       appConfig,
		authVerifier: auth.JWTVerifier{},
	}
}
