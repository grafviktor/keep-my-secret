package middleware

import (
	"net/http"

	"github.com/grafviktor/keep-my-secret/internal/api/auth"
	"github.com/grafviktor/keep-my-secret/internal/config"
)

// TokenVerifier is the interface to describe common logic. Used for dependency injection when testing the application
type TokenVerifier interface {
	VerifyAuthHeader(config config.AppConfig, w http.ResponseWriter, r *http.Request) (string, *auth.Claims, error)
}

type middleware struct {
	config       config.AppConfig
	authVerifier TokenVerifier
}

// New creates a new middleware instance
// appConfig is the application configuration
// authVerifier is the auth verifier instance which can be substituted with a mock object
// Returns new middleware instance
func New(appConfig config.AppConfig) middleware {
	return middleware{
		config:       appConfig,
		authVerifier: auth.JWTVerifier{},
	}
}
