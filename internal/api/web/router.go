package web

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	kmsMiddleware "github.com/grafviktor/keep-my-secret/internal/api/web/middleware"

	"github.com/grafviktor/keep-my-secret/internal/config"
	"github.com/grafviktor/keep-my-secret/internal/storage"
)

func NewHTTPRouter(appConfig config.AppConfig, storage storage.Storage) *chi.Mux {
	router := chi.NewRouter()
	if true {
		router.Use(chiMiddleware.Recoverer)
		router.Use(chiMiddleware.Logger)
	}

	m := kmsMiddleware.New(appConfig)

	router.Route("/api/v1", func(apiRouter chi.Router) {
		apiRouter.Use(m.EnableCORS)

		apiRouter.Route("/user", func(userRouter chi.Router) {
			apiHandler := newUserHandlerProvider(appConfig, storage)

			userRouter.Post("/register", apiHandler.RegisterHandler)
			userRouter.Post("/login", apiHandler.LoginHandler)
			userRouter.Post("/logout", apiHandler.LogoutHandler)
			userRouter.Get("/token-refresh", apiHandler.RefreshTokenHandler)
		})

		apiRouter.Route("/secrets", func(secretsRouter chi.Router) {
			secretsRouter.Use(m.AuthRequired)
			apiHandler := newSecretHandlerProvider(appConfig, storage)

			secretsRouter.Get("/", apiHandler.ListSecretsHandler)
			secretsRouter.Post("/", apiHandler.SaveSecretHandler)
			secretsRouter.Put("/{id}", apiHandler.SaveSecretHandler)
			secretsRouter.Delete("/{id}", apiHandler.DeleteSecretHandler)
			secretsRouter.Get("/file/{id}", apiHandler.DownloadSecretFileHandler)
		})

		apiRouter.Get("/version", VersionHandler)
	})

	registerStaticHandler(appConfig, router)

	return router
}
