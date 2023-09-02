package web

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/keycache"

	"github.com/golang-jwt/jwt/v4"

	"github.com/grafviktor/keep-my-secret/internal/api/auth"
	"github.com/grafviktor/keep-my-secret/internal/api/utils"
	"github.com/grafviktor/keep-my-secret/internal/config"
	"github.com/grafviktor/keep-my-secret/internal/constant"
	"github.com/grafviktor/keep-my-secret/internal/model"
)

type userStorage interface {
	AddUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, login string) (*model.User, error)
}

type userHTTPHandler struct {
	config  config.AppConfig
	storage userStorage
}

// NewApiHandler - self-explanatory
func newUserHandlerProvider(appConfig config.AppConfig, storage userStorage) userHTTPHandler {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return userHTTPHandler{
		config:  appConfig,
		storage: storage,
	}
}

type credentials struct {
	Login    string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (h *userHTTPHandler) handleSuccessFullUserSignIn(w http.ResponseWriter, user *model.User, cred credentials) {
	secret, err := user.GetDataKey(cred.Password)
	if err != nil {
		log.Printf("RegisterHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
			Status:  constant.APIStatusError,
			Message: constant.APIMessageServerError,
			Data:    nil,
		})

		return
	}

	kc := keycache.GetInstance()
	kc.Set(cred.Login, secret)

	jwtUser := auth.JWTUser{ID: cred.Login}
	authUtils := auth.New(h.config)
	tokens, err := authUtils.GenerateTokenPair(&jwtUser)
	if err != nil {
		log.Printf("LoginHandler error: cannot generate tokens. Error: %s", err.Error())

		_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
			Status:  constant.APIStatusError,
			Message: constant.APIMessageServerError,
			Data:    nil,
		})

		return
	}
	refreshCookie := authUtils.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	log.Printf("LoginUser success: Login '%s'\n", cred.Login)

	_ = utils.WriteJSON(w, http.StatusCreated, api.Response{
		Status: constant.APIStatusSuccess,
		Data:   tokens.AccessToken,
	})
}

func (h *userHTTPHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var cred credentials
	if err := utils.ReadJSON(w, r, &cred); err != nil {
		if err != nil {
			log.Printf("RegisterHandler error: %s\n", err.Error())
		}

		_ = utils.WriteJSON(w, http.StatusBadRequest, api.Response{
			Status:  constant.APIStatusFail,
			Message: constant.APIMessageBadRequest,
			Data:    nil,
		})

		return
	}

	if !utils.IsUsernameConformsPolicy(cred.Login) || !utils.IsPasswordConformsPolicy(cred.Password) {
		_ = utils.WriteJSON(w, http.StatusNotAcceptable, api.Response{
			Status:  constant.APIStatusFail,
			Message: constant.APIMessageBadRequest,
			Data:    nil,
		})

		return
	}

	user, err := model.NewUser(cred.Login, cred.Password)
	if err != nil {
		log.Printf("RegisterHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
			Status:  constant.APIStatusError,
			Message: constant.APIMessageServerError,
			Data:    nil,
		})

		return
	}

	_, err = h.storage.AddUser(r.Context(), user)
	if err != nil {
		log.Printf("RegisterHandler error: %s\n", err.Error())

		if errors.Is(err, constant.ErrDuplicateRecord) {
			_ = utils.WriteJSON(w, http.StatusConflict, api.Response{
				Status:  constant.APIStatusFail,
				Message: "user already exists",
				Data:    nil,
			})
		} else {
			_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
				Status:  constant.APIStatusError,
				Message: constant.APIMessageServerError,
				Data:    nil,
			})
		}

		return
	}

	h.handleSuccessFullUserSignIn(w, user, cred)
}

func (h *userHTTPHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var cred credentials
	// Read login and password
	if err := utils.ReadJSON(w, r, &cred); err != nil {
		_ = utils.WriteJSON(w, http.StatusBadRequest, api.Response{
			Status:  constant.APIStatusFail,
			Message: constant.APIMessageBadRequest,
			Data:    nil,
		})

		log.Printf("LoginHandler error: %s\n", err.Error())

		return
	}

	// Get password hash from the database
	user, err := h.storage.GetUser(r.Context(), cred.Login)
	if err != nil {
		log.Printf("LoginHandler error: %s\n", err.Error())

		if errors.Is(err, constant.ErrNotFound) {
			_ = utils.WriteJSON(w, http.StatusUnauthorized, api.Response{
				Status:  constant.APIStatusFail,
				Message: constant.APIMessageUnauthorized,
				Data:    nil,
			})
		} else {
			_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
				Status:  constant.APIStatusError,
				Message: constant.APIMessageServerError,
				Data:    nil,
			})
		}

		return
	}

	isPasswordCorrect, err := user.PasswordMatches(cred.Password)
	if err != nil {
		log.Printf("LoginHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
			Status:  constant.APIStatusError,
			Message: constant.APIMessageServerError,
			Data:    nil,
		})

		return
	}

	if !isPasswordCorrect {
		log.Printf("LoginHandler error: Login '%s' provided incorrect password\n", cred.Login)

		_ = utils.WriteJSON(w, http.StatusUnauthorized, api.Response{
			Status:  constant.APIStatusFail,
			Message: constant.APIMessageUnauthorized,
			Data:    nil,
		})

		return
	}

	h.handleSuccessFullUserSignIn(w, user, cred)
}

func (h *userHTTPHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Client does not provide cookie
	for _, cookie := range r.Cookies() {
		if cookie.Name == auth.CookieName {
			claims := &auth.Claims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
				return []byte(h.config.Secret), nil
			})
			if err != nil {
				log.Printf("RefreshTokenHandler error: cannot parse refresh token claims. Error: %s", err.Error())

				_ = utils.WriteJSON(w, http.StatusUnauthorized, api.Response{
					Status:  constant.APIStatusFail,
					Message: constant.APIMessageUnauthorized,
					Data:    nil,
				})

				return
			}

			_, err = h.storage.GetUser(r.Context(), claims.Subject)
			if err != nil && !errors.Is(err, constant.ErrNotFound) {
				log.Println("RefreshTokenHandler error: unknown credentials")

				_ = utils.WriteJSON(w, http.StatusUnauthorized, api.Response{
					Status:  constant.APIStatusFail,
					Message: constant.APIMessageUnauthorized,
					Data:    nil,
				})

				return
			}

			jwtUser := auth.JWTUser{ID: claims.Subject}
			a := auth.New(h.config)
			tokens, err := a.GenerateTokenPair(&jwtUser)
			if err != nil {
				log.Println("RefreshTokenHandler error: cannot create tokens")

				_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
					Status:  constant.APIStatusError,
					Message: constant.APIMessageServerError,
					Data:    nil,
				})

				return
			}

			refreshCookie := a.GetRefreshCookie(tokens.RefreshToken)
			http.SetCookie(w, refreshCookie)

			_ = utils.WriteJSON(w, http.StatusOK, api.Response{
				Status: constant.APIStatusSuccess,
				Data:   tokens.AccessToken,
			})
		}
	}
}

func (h *userHTTPHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	a := auth.New(h.config)
	refreshCookie := a.GetExpiredRefreshCookie()

	http.SetCookie(w, refreshCookie)
	_ = utils.WriteJSON(w, http.StatusOK, api.Response{
		Status: constant.APIStatusSuccess,
		Data:   nil,
	})
}
