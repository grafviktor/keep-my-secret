package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/api/utils"
	"github.com/grafviktor/keep-my-secret/internal/constant"
)

func (m middleware) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := m.authVerifier.VerifyAuthHeader(m.config, w, r)
		if err != nil {
			log.Printf("auth: %v", err.Error())

			_ = utils.WriteJSON(w, http.StatusUnauthorized, api.Response{
				Status:  constant.APIStatusFail,
				Message: "unauthorized",
				Data:    nil,
			})

			return
		}

		r = r.WithContext(context.WithValue(r.Context(), api.ContextUserLogin, claims.Subject))

		next.ServeHTTP(w, r)
	})
}
