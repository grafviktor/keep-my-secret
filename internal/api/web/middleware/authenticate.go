package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/api/utils"
	"github.com/grafviktor/keep-my-secret/internal/constant"
)

// AuthRequired middleware for checking if user is authenticated
// If user is not authenticated, it will return unauthorized response (401)
// If user is authenticated, it will add user login to context
func (m *middleware) AuthRequired(next http.Handler) http.Handler {
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
