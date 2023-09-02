package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/api/utils"
	"github.com/grafviktor/keep-my-secret/internal/config"
	"github.com/grafviktor/keep-my-secret/internal/constant"
	"github.com/grafviktor/keep-my-secret/internal/keycache"
	"github.com/grafviktor/keep-my-secret/internal/model"
	"github.com/grafviktor/keep-my-secret/internal/storage"
)

type apiRouteProvider struct {
	config  config.AppConfig
	storage storage.Storage
}

// NewApiHandler - self-explanatory
func newSecretHandlerProvider(appConfig config.AppConfig, appStorage storage.Storage) apiRouteProvider {
	return apiRouteProvider{
		config:  appConfig,
		storage: appStorage,
	}
}

func parseMultiPartSecretRequest(r *http.Request, secret *model.Secret) error {
	err := r.ParseMultipartForm(maxFileSize) // Max memory to use for parsing, in this case 10MB
	if err != nil {
		return fmt.Errorf("SaveSecretHandler error: %s", err.Error())
	}

	jsonData := r.FormValue("data")
	err = json.Unmarshal([]byte(jsonData), &secret)
	if err != nil {
		return fmt.Errorf("SaveSecretHandler error: %s", err.Error())
	}

	file, _, err := r.FormFile("file") // "file" should match the name attribute of the file input in the form
	if err != nil {
		return fmt.Errorf("SaveSecretHandler error: %s", err.Error())
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("SaveSecretHandler error: %s", err.Error())
	}

	secret.File = fileContent

	return nil
}

const maxFileSize = 1024 * 1024 * 1 // 1MB
func (a *apiRouteProvider) SaveSecretHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var secret model.Secret
	var err error

	if strings.Contains(contentType, "multipart/form-data") {
		err = parseMultiPartSecretRequest(r, &secret)
	} else {
		err = utils.ReadJSON(w, r, &secret)
	}

	if err != nil {
		log.Printf("SaveSecretHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusBadRequest, api.Response{
			Status:  constant.APIStatusFail,
			Message: constant.APIMessageBadRequest,
			Data:    nil,
		})

		return
	}

	login := r.Context().Value(api.ContextUserLogin).(string)
	key, err := keycache.GetInstance().Get(login)
	if err != nil {
		log.Printf("SaveSecretHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusUnauthorized, api.Response{
			Status:  constant.APIStatusFail,
			Message: constant.APIMessageUnauthorized,
			Data:    nil,
		})

		return
	}

	err = secret.Encrypt(key, login+a.config.Secret)
	if err != nil {
		log.Printf("SaveSecretHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
			Status:  constant.APIStatusError,
			Message: constant.APIMessageServerError,
			Data:    nil,
		})

		return
	}

	_, err = a.storage.SaveSecret(r.Context(), &secret, login)

	if err != nil {
		log.Printf("SaveSecretHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
			Status:  constant.APIStatusError,
			Message: constant.APIMessageServerError,
			Data:    nil,
		})

		return
	}

	_ = utils.WriteJSON(w, http.StatusCreated, api.Response{
		Status: constant.APIStatusSuccess,
		Data:   secret,
	})
}

func (a *apiRouteProvider) ListSecretsHandler(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(api.ContextUserLogin).(string)
	key, err := keycache.GetInstance().Get(login)
	if err != nil {
		log.Printf("ListSecretsHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusUnauthorized, api.Response{
			Status:  constant.APIStatusFail,
			Message: constant.APIMessageUnauthorized,
			Data:    nil,
		})

		return
	}

	secrets, err := a.storage.GetSecretsByUser(r.Context(), login)
	if err != nil {
		log.Printf("ListSecretsHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
			Status:  constant.APIStatusError,
			Message: constant.APIMessageServerError,
			Data:    nil,
		})

		return
	}

	for _, secret := range secrets {
		err = secret.Decrypt(key, login+a.config.Secret)
		if err != nil {
			log.Printf("ListSecretsHandler error: %s\n", err.Error())

			_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
				Status:  constant.APIStatusError,
				Message: constant.APIMessageServerError,
				Data:    nil,
			})

			return
		}
	}

	_ = utils.WriteJSON(w, http.StatusOK, api.Response{
		Status: constant.APIStatusSuccess,
		Data:   secrets,
	})
}

func (a *apiRouteProvider) DeleteSecretHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	login := r.Context().Value(api.ContextUserLogin).(string)

	err := a.storage.DeleteSecret(r.Context(), id, login)
	if err != nil {
		log.Printf("DeleteSecretHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusNotFound, api.Response{
			Status:  constant.APIStatusFail,
			Message: constant.APIMessageNotFound,
			Data:    nil,
		})

		return
	}

	_ = utils.WriteJSON(w, http.StatusAccepted, api.Response{
		Status: constant.APIStatusSuccess,
		Data:   id,
	})
}

func (a *apiRouteProvider) DownloadSecretFileHandler(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value(api.ContextUserLogin).(string)
	secretID := chi.URLParam(r, "id")

	key, err := keycache.GetInstance().Get(login)
	if err != nil {
		log.Printf("DownloadSecretFileHandler error: %s\n", err.Error())

		http.Error(w, constant.APIMessageUnauthorized, http.StatusUnauthorized)

		return
	}

	secret, err := a.storage.GetSecret(r.Context(), secretID, login)
	if err != nil {
		log.Printf("DownloadSecretFileHandler error: %s\n", err.Error())

		// This is not a JSON content-type handler
		if errors.Is(err, constant.ErrNotFound) {
			http.Error(w, constant.APIMessageNotFound, http.StatusNotFound)
		} else {
			http.Error(w, constant.APIMessageServerError, http.StatusInternalServerError)
		}

		return
	}

	err = secret.Decrypt(key, login+a.config.Secret)
	if err != nil {
		log.Printf("DownloadSecretFileHandler error: %s\n", err.Error())

		_ = utils.WriteJSON(w, http.StatusInternalServerError, api.Response{
			Status:  constant.APIStatusError,
			Message: constant.APIMessageServerError,
			Data:    nil,
		})

		return
	}

	// Set headers for the download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", secret.FileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(secret.File)))

	// Stream the file content to the response
	_, err = io.Copy(w, bytes.NewReader(secret.File))
	if err != nil {
		http.Error(w, "Error streaming file content to response", http.StatusInternalServerError)
		return
	}
}
