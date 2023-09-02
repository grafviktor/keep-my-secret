package middleware

import "github.com/grafviktor/keep-my-secret/internal/config"

type middleware struct {
	config config.AppConfig
}

func New(appConfig config.AppConfig) middleware {
	return middleware{
		config: appConfig,
	}
}
